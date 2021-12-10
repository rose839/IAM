package secret

import (
	"github.com/gin-gonic/gin"
	metav1 "github.com/rose839/IAM/api/meta/v1"
	"github.com/rose839/IAM/internal/pkg/middleware"
	"github.com/rose839/IAM/pkg/core"
)

// DeleteCollection delete secrets by secret names.
func (s *SecretController) DeleteCollection(c *gin.Context) {

	if err := s.srv.Policies().DeleteCollection(
		c,
		c.GetString(middleware.UsernameKey),
		c.QueryArray("name"),
		metav1.DeleteOptions{},
	); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
