package middleware

import "github.com/gin-gonic/gin"

const (
	RequestIDKey = "requestID"
	UsernameKey  = "username"
)

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(RequestIDKey, c.GetString(XRequestIDKey))
		c.Set(UsernameKey, c.GetString(UsernameKey))
		c.Next()
	}
}
