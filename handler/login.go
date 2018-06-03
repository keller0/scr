package handle

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var err error
	var loginJSON login
	if err = c.ShouldBindJSON(&loginJSON); err == nil {
		fmt.Println(loginJSON)

		Con, err := sql.Open("mysql", "root:111@tcp(127.0.0.1:3306)/yxi")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer Con.Close()

		var password string
		var id int64
		err = Con.QueryRow("SELECT id, password FROM user WHERE username=?", loginJSON.User).Scan(
			&id, &password)
		if err != nil {
			fmt.Println(err.Error())
		}

		if loginJSON.Password == password {
			// password is right, generate token

			token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
			// Set some claims
			token.Claims = jwt_lib.MapClaims{
				"Id":       id,
				"username": loginJSON.User,
				"exp":      time.Now().Add(time.Hour * 1).Unix(),
			}
			// Sign and get the complete encoded token as a string
			tokenString, err := token.SignedString([]byte("secret"))
			if err != nil {
				c.JSON(500, gin.H{"message": "Could not generate token"})
			}
			c.JSON(http.StatusOK, gin.H{"token": tokenString})

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
