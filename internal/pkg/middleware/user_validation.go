package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "github.com/rose839/IAM/api/meta/v1"
	"github.com/rose839/IAM/internal/apiserver/store"
	"github.com/rose839/IAM/internal/pkg/code"
	"github.com/rose839/IAM/pkg/core"
	"github.com/rose839/IAM/pkg/errors"
)

// Validation make sure users have the right resource permission and operation.
func Validation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// normal users, not admin
		if err := isAdmin(c); err != nil {
			switch c.FullPath() {
			case "/v1/users":
				if c.Request.Method != http.MethodPost {
					core.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, ""), nil)
					c.Abort()

					return
				}
			case "/v1/users/:name", "/v1/users/:name/change_password":
				username := c.GetString("username")
				if c.Request.Method == http.MethodDelete ||
					(c.Request.Method != http.MethodDelete && username != c.Param("name")) {
					core.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, ""), nil)
					c.Abort()

					return
				}
			default:
			}
		}

		c.Next()
	}
}

// isAdmin make sure the user is administrator.
func isAdmin(c *gin.Context) error {
	username := c.GetString(UsernameKey)
	user, err := store.Client().Users().Get(c, username, metav1.GetOptions{})
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	if user.IsAdmin != 1 {
		return errors.WithCode(code.ErrPermissionDenied, "user %s is not a administrator", username)
	}

	return nil
}
