package util

import (
	"errors"
	"os"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	authSecret = os.Getenv("YXI_BACK_KEY")
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
