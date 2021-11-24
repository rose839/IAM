package v1

import (
	"encoding/json"

	metav1 "github.com/rose839/IAM/api/meta/v1"
	"github.com/rose839/IAM/pkg/auth"
	"gorm.io/gorm"
)

// User represents a user restful resource. It is also used as gorm model.
type User struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Required: true
	Nickname string `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`

	// Required: true
	Password string `json:"password,omitempty" gorm:"column:password" validate:"required"`

	// Required: true
	Email string `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`

	Phone string `json:"phone" gorm:"column:phone" validate:"omitempty"`

	IsAdmin int `json:"isAdmin,omitempty" gorm:"column:isAdmin" validate:"omitempty"`

	TotalPolicy int64 `json:"totalPolicy" gorm:"-" validate:"omitempty"`
}

// UserList is the whole list of all users which have been stored in stroage.
type UserList struct {
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*User `json:"iterms"`
}

// TableName maps to mysql table name.
func (u *User) TableName() string {
	return "user"
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

// BeforeCreate run before create database record.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	u.ExtendShadow = u.Extend.String()
	return
}

// AfterCreate run after create database record.
func (u *User) AfterCreate(tx *gorm.DB) (err error) {

}

// BeforeUpdate run before update database record.
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	u.ExtendShadow = u.Extend.String()

	return
}

// AfterFind run after find to unmarshal a extend shadown string into metav1.Extend struct.
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	if err := json.Unmarshal([]byte(u.ExtendShadow), &u.Extend); err != nil {
		return err
	}

	return nil
}
