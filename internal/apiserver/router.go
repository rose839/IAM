package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/rose839/IAM/internal/apiserver/controller/v1/policy"
	"github.com/rose839/IAM/internal/apiserver/controller/v1/secret"
	"github.com/rose839/IAM/internal/apiserver/controller/v1/user"
	"github.com/rose839/IAM/internal/apiserver/store/mysql"
	"github.com/rose839/IAM/internal/pkg/code"
	"github.com/rose839/IAM/internal/pkg/middleware"
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

func installController(g *gin.Engine) *gin.Engine {
	JWTStrategy, _ := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", JWTStrategy.LoginHandler)
	g.POST("logout", JWTStrategy.LogoutHandler)
	g.POST("/refresh", JWTStrategy.RefreshHandler) // Refresh time can be longer than token timeout

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	// v1 handlers, requiring authentication
	storeIns, _ := mysql.GetMySQLFactoryOr(nil)
	v1 := g.Group("/v1")
	v1.Use(auto.AuthFunc())
	{
		// user RESTful resource
		userv1 := v1.Group("/users")
		userv1.Use(auto.AuthFunc(), middleware.Validation())
		{
			userController := user.NewUserController(storeIns)

			userv1.POST("", userController.Create)
			userv1.DELETE("", userController.Delete)      // admin api
			userv1.DELETE(":name", userController.Delete) // admin api
			userv1.PUT(":name/change-password", userController.ChangePassword)
			userv1.PUT(":name", userController.Update)
			userv1.GET("", userController.List)
			userv1.GET(":name", userController.Get)
		}

		// police RESTful resource
		policyv1 := v1.Group("/policies", middleware.Publish())
		{
			policyController := policy.NewPolicyController(storeIns)

			policyv1.POST("", policyController.Create)
			policyv1.DELETE("", policyController.Delete)
			policyv1.DELETE(":name", policyController.Delete)
			policyv1.PUT(":name", policyController.Update)
			policyv1.GET("", policyController.List)
			policyv1.GET(":name", policyController.Get)
		}

		// secret RESTful resource
		secretv1 := v1.Group("/secrets", middleware.Publish())
		{
			secretController := secret.NewSecretController(storeIns)

			secretv1.POST("", secretController.Create)
			secretv1.DELETE(":name", secretController.Delete)
			secretv1.PUT(":name", secretController.Update)
			secretv1.GET("", secretController.List)
			secretv1.GET(":name", secretController.Get)
		}
	}

	return g
}
