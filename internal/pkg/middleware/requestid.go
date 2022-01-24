package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	// XRequestIDKey defines X-Request-ID key string.
	XRequestIDKey = "X-Request-ID"
)

// Add request id if no request id found in request headers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		rid := c.GetHeader(XRequestIDKey)

		if rid == "" {
			rid = uuid.Must(uuid.NewV4(), nil).String()
			c.Request.Header.Set(XRequestIDKey, rid)
			c.Set(XRequestIDKey, rid)
		}

		// Set XRequestIDKey response header
		c.Writer.Header().Set(XRequestIDKey, rid)
		c.Next()
	}
}
