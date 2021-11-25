package v1

import (
	"github.com/rose839/IAM/pkg/validation/"
	"github.com/rose839/IAM/pkg/validation/field"
)

func (u *User) Validate() field.ErrorList {
	val := validation.NewValidator(u)
	allErrs := val.Validate()

	if err := validation.IsValidPassword(u.Password); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("password"), err.Error(), ""))
	}

	return allErrs
}
