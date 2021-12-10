package secret

import (
	srvv1 "github.com/rose839/IAM/internal/apiserver/service/v1"
	"github.com/rose839/IAM/internal/apiserver/store"
)

// SecretController create a secret handler used to handle request for secret resource.
type SecretController struct {
	srv srvv1.Service
}

// NewSecretController creates a secret handler.
func NewSecretController(store store.Factory) *SecretController {
	return &SecretController{
		srv: srvv1.NewService(store),
	}
}
