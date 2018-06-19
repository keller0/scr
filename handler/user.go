package handle

import (
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/model"
)

type login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type register struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Login return a jwt if user info is valid.
// 200 400 401 404
func Login(c *gin.Context) {
	var err error
	var loginJSON login
	if err = c.ShouldBindJSON(&loginJSON); err == nil {
		var user model.User
		user.Username = loginJSON.User
		user.Password = loginJSON.Password

		if !user.UsernameExist() {
			// return if username allready exists
			c.JSON(http.StatusNotFound, gin.H{"error": "user dose not exists"})
			c.Abort()
			return
		}
		tokenString, err := user.Login()
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

// CheckUserExist use http post check if username already exists
func CheckUserExist(c *gin.Context) {
	var username string
	if err := c.ShouldBindQuery(username); err != nil {
		var user model.User
		user.Username = username
		if user.UsernameExist() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": ""})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// Register use post data to create a user account
// 200 400 409 500
func Register(c *gin.Context) {
	var err error
	var registJSON register
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	reName := regexp.MustCompile("^[a-zA-Z0-9]+$")
	err = c.ShouldBindJSON(&registJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if !re.MatchString(registJSON.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not valid."})
		c.Abort()
		return
	}
	if !reName.MatchString(registJSON.User) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is not valid."})
		c.Abort()
		return
	}
	var user model.User
	user.Username = registJSON.User
	user.Email = registJSON.Email

	if user.UsernameExist() {
		// return if username allready exists
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		c.Abort()
		return
	}
	if user.EmailExist() {
		// return if username allready exists
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		c.Abort()
		return
	}

	user.Password = registJSON.Password
	e := user.New()
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
	}
	c.String(http.StatusOK, "registration succeeded")
}
