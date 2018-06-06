package util

import (
	"errors"
	"math/rand"
	"net/http"
	"os"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"golang.org/x/crypto/bcrypt"
)

var (
	authSecret = os.Getenv("YXI_BACK_KEY")

	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()")
)

// CheckPasswordHash check if passwoed match bcrypt hash.
// return true if match.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JwtGenToken gnerate new token add claims
// return signed string and nil if succeed.
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
		return -1, err
	}

	if claims, ok := token.Claims.(jwt_lib.MapClaims); ok && token.Valid {
		id := claims["id"].(float64)
		//got userid
		return int64(id), nil
	}

	return -1, errors.New("get id failed")
}

// HashPassword use bcrypt hash user's password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	return string(bytes), err
}

// RandStringRunes return n bits rand string
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
