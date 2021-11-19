package apiserver

import "github.com/gin-gonic/gin"

// init rest api server router.
func initRouter(g *gin.Engine) {
	InstallMiddleware(g)
	installController(g)
}

func InstallMiddleware(g *gin.Engine) {

}

func installController(g *gin.Engine) {

}
