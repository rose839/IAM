package user

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/rose839/IAM/api/apiserver/v1"
	"github.com/rose839/IAM/internal/pkg/code"
	"github.com/rose839/IAM/pkg/core"
	"github.com/rose839/IAM/pkg/errors"
)

func (u *UserController) Create(c *gin.Context) {
	var r v1.User

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
	}
}
