package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ErrorLimitExceeded defines Limit exceeded error.
var ErrorLimitExceeded = errors.New("Limit exceeded")

//
func Limit(maxEventsPerSec float64, maxBurstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(maxEventsPerSec), maxBurstSize)

	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()

			return
		}

		// Limit reached
		_ = c.Error(ErrorLimitExceeded)
		c.AbortWithStatus(429)
	}
}
