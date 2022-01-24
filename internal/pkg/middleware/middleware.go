package middleware

import (
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

// Middlewares store registered middlewares.
var Middlewares = defaultMiddlewares()

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":  gin.Recovery(),
		"cors":      Cors(),
		"requestid": RequestID(),
		"dump":      gindump.Dump(),
	}
}
