package apiserver

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/rose839/IAM/internal/pkg/middleware"
	"github.com/rose839/IAM/internal/pkg/middleware/auth"
	"github.com/spf13/viper"
)

const (
	// APIServerAudience defines the value of jwt audience field.
	APIServerAudience = "iam.api.rose839.com"

	// APIServerIssuer defines the value of jwt issuer field.
	APIServerIssuer = "iam-apiserver"
)

type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username string, password string) bool {

	})
}

func newJWTAuth() middleware.AuthStrategy {
	ginjwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm: viper.GetString("jwt.Realm")
		SigningAlgorithm: "HS256",
	})

	return auth.NewJWTStrategy(*ginjwt)
}
