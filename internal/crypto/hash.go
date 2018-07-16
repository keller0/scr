package crypto

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CheckPasswordHash check if passwoed match bcrypt hash.
// return true if match.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPassword use bcrypt hash user's password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	return string(bytes), err
}

// RandString return n bits rand string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
