package secret

import (
	"github.com/gin-gonic/gin"
	metav1 "github.com/rose839/IAM/api/meta/v1"
	"github.com/rose839/IAM/internal/pkg/middleware"
	"github.com/rose839/IAM/pkg/core"
)

// Get get an policy by the secret identifier.
func (s *SecretController) Get(c *gin.Context) {
	secret, err := s.srv.Secrets().Get(c, c.GetString(middleware.UsernameKey), c.Param("name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, secret)
}
