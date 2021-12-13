package v1

import (
	"github.com/rose839/IAM/pkg/validation"
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

// ValidateUpdate validates that a user object is valid when update.
// Like User.Validate but not validate password.
func (u *User) ValidateUpdate() field.ErrorList {
	val := validation.NewValidator(u)
	allErrs := val.Validate()

	return allErrs
}

func (s *Secret) Validate() field.ErrorList {
	val := validation.NewValidator(s)

	return val.Validate()
}

// Validate validates that a policy object is valid.
func (p *Policy) Validate() field.ErrorList {
	val := validation.NewValidator(p)

	return val.Validate()
}
