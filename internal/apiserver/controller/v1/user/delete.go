package user

import (
	"github.com/gin-gonic/gin"
	metav1 "github.com/rose839/IAM/api/meta/v1"
	"github.com/rose839/IAM/pkg/core"
)

// Delete delete an user by the user identifier.
// Only administrator can call this function.
func (u *UserController) Delete(c *gin.Context) {
	if err := u.srv.Users().Delete(c, c.Param("name"), metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
