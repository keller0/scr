package token

import (
	"errors"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/keller0/yxi-back/internal"
	"net/http"
)

var (
	authSecret = internal.GetEnv("YXI_BACK_KEY", "secretkey")
)

// JwtGenToken gnerate new token add claims
// return signed string and nil if succeed.
func JwtGenToken(userID int64, userName, runToken string, exp int64) (string, error) {
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	// Set some claims
	token.Claims = jwt_lib.MapClaims{
		"id":       userID,
		"username": userName,
		"runtoken": runToken,
		"exp":      exp,
	}
	// Sign and get the complete encoded token as a string
	return token.SignedString([]byte(authSecret))
}

// JwtOK check request's jwt is valid,
func JwtOK(r *http.Request) (bool, error) {
	_, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
		b := ([]byte(authSecret))
		return b, nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// JwtGetUserID return userid from jwt claims
func JwtGetUserID(r *http.Request) (int64, error) {

	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
		b := ([]byte(authSecret))
		return b, nil
	})
	if err != nil {
		// get token failed
		return -1, err
	}

	if claims, ok := token.Claims.(jwt_lib.MapClaims); ok && token.Valid {
		id := claims["id"].(float64)
		//got userid
		return int64(id), nil
	}

	return -1, errors.New("get id failed")
}
