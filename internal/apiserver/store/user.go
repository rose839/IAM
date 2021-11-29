package store

import (
	"context"

	v1 "github.com/rose839/IAM/api/apiserver/v1"
	metav1 "github.com/rose839/IAM/api/meta/v1"
)

// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
	Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1.GetOptions) error
	List(ctx context.Context, opts metav1.ListOption) error
}
