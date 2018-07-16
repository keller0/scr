package model

import (
	"log"
	"time"

	"github.com/keller0/yxi-back/db/mysql"
	"github.com/keller0/yxi-back/db/redis"
	"github.com/keller0/yxi-back/internal/crypto"
	"github.com/keller0/yxi-back/internal/token"
	"github.com/keller0/yxi-back/service/mailgun"
)

var (
	redisPrefixResetPassEmail = "reset-pass-"
)

// User user struct in database
type User struct {
	ID int64 `json:"id"`

	Username string `json:"username"`

	Password string `json:"password"`

	RunToke string `json:"run_token"`

	Email string `json:"email"`

	CreateAt string `json:"createat"`

	UpdateAt string `json:"updateat"`

	Admin bool `json:"admin"`
}

// UsernameExist check if username already existed
func (u *User) UsernameExist() bool {

	var count int64
	err := mysql.Db.QueryRow("SELECT count(id) FROM user WHERE username=?", u.Username).Scan(&count)
	if err != nil {
		return true
	}
	return count != 0
}

// EmailExist check if email already existed
func (u *User) EmailExist() bool {
	return emailExist(u.Email)
}

func emailExist(email string) bool {

	var count int64
	err := mysql.Db.QueryRow("SELECT count(id) FROM user WHERE email=?", email).Scan(&count)
	if err != nil {
		return true
	}
	return count != 0
}

// Login check user password and return jwt
func (u *User) Login() (string, error) {

	var password, runToken string
	var id int64
	var err error
	err = mysql.Db.QueryRow("SELECT id, password, run_token FROM user WHERE username=?",
		u.Username).Scan(&id, &password, &runToken)
	if err != nil {
		return "", err
	}

	if crypto.CheckPasswordHash(u.Password, password) {
		exp := time.Now().Add(time.Hour * 24 * 15).Unix()
		tokenString, err := token.JwtGenToken(id, u.Username, runToken, exp)
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}

	return "", ErrLoginWrongPass
}

// New create a new user account
func (u *User) New() error {

	var runToken = crypto.RandString(40)
	passwordhashed, err := crypto.HashPassword(u.Password)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	insUser, err := mysql.Db.Prepare("INSERT INTO user(username, password, email, run_token) values(?,?,?,?)")
	defer insUser.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, e := insUser.Exec(u.Username, passwordhashed, u.Email, runToken)
	return e
}

// SendResetToken generate a token, store it in redis and send msg to email
func SendResetToken(email string) (err error) {
	if emailExist(email) != true {
		return ErrEmailNotExist
	}

	token := crypto.RandString(60)
	err = redis.Set(redisPrefixResetPassEmail+email, []byte(token))
	if err != nil {
		return
	}
	// set expire = 2 days
	err = redis.Expire(redisPrefixResetPassEmail+email, 3600*24*2)
	if err != nil {
		return
	}

	msg := "To reset your password please click this link https://yxi.io/account?email=" + email + "&token=" + token
	id, err := mailgun.SimpleMessage("Reset password", msg, email)
	if err != nil {
		return
	}
	log.Println("send-mail id :", id)
	return
}
