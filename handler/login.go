package handle

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/middleware"
)

type login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var err error
	var loginJSON login
	if err = c.ShouldBindJSON(&loginJSON); err == nil {

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
			log.Fatal(err.Error())
		}

		if loginJSON.Password == password {
			exp := time.Now().Add(time.Hour * 1).Unix()
			tokenString, err := mid.JwtGenToken(id, loginJSON.User, exp)
			if err != nil {
				c.JSON(500, gin.H{"message": "Could not generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": tokenString})

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
