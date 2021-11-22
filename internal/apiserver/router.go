package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/rose839/IAM/internal/pkg/code"
	"github.com/rose839/IAM/internal/pkg/middleware/auth"
	"github.com/rose839/IAM/pkg/core"
	"github.com/rose839/IAM/pkg/errors"
)

// init rest api server router.
func initRouter(g *gin.Engine) {
	InstallMiddleware(g)
	installController(g)
}

func InstallMiddleware(g *gin.Engine) {

}

func installController(g *gin.Engine) {
	JWTStrategy, _ := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", JWTStrategy.LoginHandler)
	g.POST("logout", JWTStrategy.LogoutHandler)
	g.POST("/refresh", JWTStrategy.RefreshHandler) // Refresh time can be longer than token timeout

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	// v1 handlers, requiring authentication
	v1 := g.Group("/v1")
	v1.Use(auto.AuthFunc())
	{
		// user RESTful resource
		userv1 := v1.Group("/users")
		{
		}

		// police RESTful resource
		policyv1 := v1.Group("/policies")
		{
		}

		// secret RESTful resource
		secretv1 := v1.Group("/secrets")
		{
		}
	}
}
