package middleware

import "github.com/gin-gonic/gin"

// Publish publish a redis event to specified redis channel when some action occurred.
func Publish() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
