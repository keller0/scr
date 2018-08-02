package handle

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/model"
)

type login struct {
	//user could be username or email
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type registerMail struct {
	User  string `form:"user" json:"user" binding:"required"`
	Email string `form:"email" json:"email" binding:"required"`
}
type register struct {
	User  string `form:"user" json:"user" binding:"required"`
	Email string `form:"email" json:"email" binding:"required"`
	Pass  string `form:"pass" json:"pass" binding:"required"`
	Token string `form:"token" json:"token" binding:"required"`
}

type resetMail struct {
	Email string `form:"email" json:"email" binding:"required"`
}

type resetPass struct {
	Email string `form:"email" json:"email" binding:"required"`
	Pass  string `form:"pass" json:"pass" binding:"required"`
	Token string `form:"token" json:"token" binding:"required"`
}

// Login return a jwt if user info is valid.
// 200 400 401 404
func Login(c *gin.Context) {
	var err error
	var loginJSON login
	if err = c.ShouldBindJSON(&loginJSON); err == nil {
		var user model.User
		if strings.Index(loginJSON.User, "@") != -1 {
			user.Email = loginJSON.User
			user.Username = ""
			if !user.EmailExist() {
				// return if email not exists
				c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["EmailNotExist"]})
				return
			}
		} else {
			user.Username = loginJSON.User
			user.Email = ""
			if !user.UsernameExist() {
				// return if username not exists
				c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["UserNotExist"]})
				return
			}
		}
		user.Password = loginJSON.Password

		tokenString, err := user.Login()
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"errNumber": responseErr["Wrong Password"]})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Bad Requset"]})
	}

}

// CheckUserExist use http post check if username already exists
func CheckUserExist(c *gin.Context) {
	var username string
	if err := c.ShouldBindQuery(username); err != nil {
		var user model.User
		user.Username = username
		if user.UsernameExist() {
			c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["User Already Exist"]})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": ""})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Bad Requset"]})
	}
}

// SendRegisterEmail use post data to create a tmp user account
// send email then wait user to verify
// 200 400 409 500
func SendRegisterEmail(c *gin.Context) {
	var err error
	var registJSON registerMail

	err = c.ShouldBindJSON(&registJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Bad Requset"]})
		return
	}

	if es := registJSON.Validate(); es != "" {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": es})
		return
	}
	var user model.User
	user.Username = registJSON.User
	user.Email = registJSON.Email

	if user.UsernameExist() {
		// return if username already exists
		c.JSON(http.StatusConflict, gin.H{"errNumber": responseErr["User Already Exist"]})
		return
	}
	if user.EmailExist() {
		// return if username already exists
		c.JSON(http.StatusConflict, gin.H{"errNumber": responseErr["Email Already Exist"]})
		return
	}

	e := user.SendRegisterToken()
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["Send Register email Failed"]})
	}
	c.String(http.StatusOK, "send registration email succeeded")
}

// RegisterComplete use post data to create a user account
// 200 400 409 500
func RegisterComplete(c *gin.Context) {
	var err error
	var registerJSON register

	err = c.ShouldBindJSON(&registerJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Bad Requset"]})
		return
	}
	if es := registerJSON.Validate(); es != "" {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": es})
		return
	}
	var user model.User
	user.Username = registerJSON.User
	user.Email = registerJSON.Email
	user.Password = registerJSON.Pass

	err = user.New(registerJSON.Token)
	if err != nil {
		if err == model.ErrTokenNotMatch {
			c.JSON(http.StatusUnauthorized, gin.H{"errNumber": responseErr["RegisterTokenNotMatch"]})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["ServerErr Register Failed"]})
		}
	}
	c.String(http.StatusOK, "registration succeeded")
}

func (r *registerMail) Validate() string {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	reName := regexp.MustCompile("^[a-zA-Z0-9]+$")
	switch {
	case !re.MatchString(r.Email):
		return responseErr["Email is not valid"]
	case !reName.MatchString(r.User):
		return responseErr["Username is not valid"]
	case len(r.User) > 15:
		return responseErr["Username is too long"]
	}
	return ""
}

func (r *register) Validate() string {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	reName := regexp.MustCompile("^[a-zA-Z0-9]+$")
	switch {
	case !re.MatchString(r.Email):
		return responseErr["Email is not valid"]
	case !reName.MatchString(r.User):
		return responseErr["Username is not valid"]
	case len(r.Pass) < 9:
		return responseErr["Password is too short"]
	}
	return ""
}

// UpdatePassByEmail reset user's password use email and token
func UpdatePassByEmail(c *gin.Context) {
	var j resetPass
	err := c.ShouldBindJSON(&j)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Bad Requset"]})
		return
	}
	err = model.UpdatePassByToken(j.Email, j.Token, j.Pass)
	if err != nil {
		log.Println(err)
		if err == model.ErrTokenNotMatch {
			c.JSON(http.StatusUnauthorized, gin.H{"errNumber": responseErr["ResetTokenNotMatch"]})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["Update Password Failed"]})
		}
		return
	}
	c.String(http.StatusOK, "update password succeeded")
	return
}

// SendResetPassEmail send reset password link to email
// 400 500 200
func SendResetPassEmail(c *gin.Context) {
	var j resetMail
	err := c.ShouldBindJSON(&j)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Bad Requset"]})
		return
	}
	err = model.SendResetToken(j.Email)
	if err != nil {
		if err == model.ErrEmailNotExist {
			c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["EmailNotExist"]})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["Send reset email Failed"]})
		}
		return
	}
	c.String(http.StatusOK, "send email succeeded")
}
