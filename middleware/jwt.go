package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi.io/internal/token"
)

// JwtAuth only alow requests with jwt
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := token.JwtGetUserID(c.Request)
		if err != nil {
			c.AbortWithError(401, err)
		}
		c.Set("uid", id)
	}
}
