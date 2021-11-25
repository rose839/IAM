package validation

import (
	"fmt"
	"strings"
	"unicode"
)

func IsQualifiedName(value string) []string {
	return nil
}

const (
	minPassLength = 8
	maxPassLength = 16
)

// IsValidPassword validate password.
func IsValidPassword(password string) error {
	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			hasNumber = true
			passLen++
		case unicode.IsUpper(ch):
			hasUpper = true
			passLen++
		case unicode.IsLower(ch):
			hasLower = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}

	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !hasLower {
		appendError("lowercase letter missing")
	}
	if !hasUpper {
		appendError("uppercase letter missing")
	}
	if !hasNumber {
		appendError("at least one numeric character required")
	}
	if !hasSpecial {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(
			fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength),
		)
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}

	return nil
}
