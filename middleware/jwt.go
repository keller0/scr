package mid

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/util"
)

var (
	authSecret = os.Getenv("YXI_BACK_KEY")
)

// JwtAuth only alow requests with jwt
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := util.JwtOK(c.Request)
		if err != nil && !ok {
			c.AbortWithError(401, err)
		}
	}
}
