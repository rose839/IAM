package middleware

import (
	"github.com/gin-gonic/gin"
)

// Define Redis pub/sub events.
const (
	RedisPubSubChannel  = "iam.cluster.notifications"
	NoticePolicyChanged = "PolicyChanged"
	NoticeSecretChanged = "SecretChanged"
)

// Publish publish a redis event to specified redis channel when some action occurred.
func Publish() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func notify(method string, command string) {

}
