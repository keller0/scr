package mid

import (
	"errors"
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

func JwtGenToken(userID int64, userName string, exp int64) (string, error) {
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	// Set some claims
	token.Claims = jwt_lib.MapClaims{
		"id":       userID,
		"username": userName,
		"exp":      exp,
	}
	// Sign and get the complete encoded token as a string
	return token.SignedString([]byte(authSecret))
}

func JwtGetUserID(tokenString string) (int64, error) {

	token, err := jwt_lib.Parse(tokenString[7:],
		func(t *jwt_lib.Token) (interface{}, error) { return []byte(authSecret), nil })
	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(jwt_lib.MapClaims); ok && token.Valid {
		id := claims["id"].(float64)
		return int64(id), nil
	}
	return -1, errors.New("get id failed")
}
