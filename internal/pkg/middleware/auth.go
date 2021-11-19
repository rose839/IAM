package middleware

import "github.com/gin-gonic/gin"

// AuthStrategy defines the set of methods used to do resource authentication.
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}
