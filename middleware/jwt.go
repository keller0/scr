package mid

import (
	"os"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

var (
	authSecret = os.Getenv("YXI_BACK_KEY")
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
			b := ([]byte(authSecret))
			return b, nil
		})

		if err != nil {
			c.AbortWithError(401, err)
		}
	}
}
