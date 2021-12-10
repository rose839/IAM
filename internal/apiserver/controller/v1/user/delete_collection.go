package user

import (
	"github.com/gin-gonic/gin"
	metav1 "github.com/rose839/IAM/api/meta/v1"
	"github.com/rose839/IAM/pkg/core"
)

// DeleteCollection batch delete users by multiple usernames.
// Only administrator can call this function.
func (u *UserController) DeleteCollection(c *gin.Context) {
	usernames := c.QueryArray("name")

	if err := u.srv.Users().DeleteCollection(c, usernames, metav1.DeleteOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
