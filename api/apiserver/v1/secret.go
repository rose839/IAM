package v1

import metav1 "github.com/rose839/IAM/api/meta/v1"

// Secret represents a secret restful resource.
// It is also used as gorm model.
type Secret struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Username  string `json:"username" gorm:"column:username" validate:"omitempty"`
	SecretID  string `json:"secretID" gorm:"column:secretID" validate:"omitempty"`
	SecretKey string `json:"secretKey" gorm:"column:secretKey" validate:"omitempty"`

	// Requires: true
	Expires     int64  `json:"expires" gorm:"column:expires" validate:"omitempty"`
	Description string `json:"description" gorm:"column:description" validate:"description"`
}

// SecretList is the whole list of all secrets which have been stored in stroage.
type SecretList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	metav1.ListMeta `json:",inline"`

	// List of secrets
	Items []*Secret `json:"items"`
}

// TableName maps to mysql table name.
func (s *Secret) TableName() string {
	return "secret"
}

