package secret

import (
	"github.com/gin-gonic/gin"
	metav1 "github.com/rose839/IAM/api/meta/v1"
	"github.com/rose839/IAM/internal/pkg/middleware"
	"github.com/rose839/IAM/pkg/core"
)

// Delete delete a secret by the secret identifier.
func (s *SecretController) Delete(c *gin.Context) {
	opts := metav1.DeleteOptions{Unscoped: true}
	if err := s.srv.Secrets().Delete(c, c.GetString(middleware.UsernameKey), c.Param("name"), opts); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
