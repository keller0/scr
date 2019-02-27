package model

import (
	"log"
	"net/url"
	"time"

	"github.com/keller0/yxi.io/db/mysql"
	"github.com/keller0/yxi.io/db/redis"
	"github.com/keller0/yxi.io/internal/crypto"
	"github.com/keller0/yxi.io/internal/token"
	"github.com/keller0/yxi.io/service/mailgun"
)

var (
	redisPrefixResetPassEmail = "reset-pass-"
	redisPrefixRegister       = "register-"
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

	var password, runToken, username string
	var id int64
	var err error
	if u.Username != "" {
		// use username to login
		err = mysql.Db.QueryRow("SELECT id, username, password, run_token FROM user WHERE username=?",
			u.Username).Scan(&id, &username, &password, &runToken)
	} else {
		// use email to login
		err = mysql.Db.QueryRow("SELECT id, username, password, run_token FROM user WHERE email=?",
			u.Email).Scan(&id, &username, &password, &runToken)
	}
	if err != nil {
		return "", err
	}
	if crypto.CheckPasswordHash(u.Password, password) {
		exp := time.Now().Add(time.Hour * 24 * 15).Unix()
		tokenString, err := token.JwtGenToken(id, username, runToken, exp)
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}

	return "", ErrLoginWrongPass
}

// SendRegisterToken generate a token, store it in redis and send msg to email
func (u *User) SendRegisterToken() (err error) {

	token := crypto.RandString(30)
	err = redis.Set(redisPrefixRegister+u.Email, []byte(token))
	if err != nil {
		return
	}
	err = redis.Expire(redisPrefixRegister+u.Email, 3600*24*2)
	if err != nil {
		return
	}

	es := url.QueryEscape(u.Email)
	ts := url.QueryEscape(token)
	us := url.QueryEscape(u.Username)

	msg := "To complete your account registration please click this link " +
		"https://www.yxi.io/singup_complete/?email=" + es + "&token=" + ts + "&user=" + us
	id, err := mailgun.SimpleMessage("Complete your account registration", msg, u.Email)
	if err != nil {
		return
	}
	log.Println("send-mail id :", id)
	return
}

// New create a new user account
func (u *User) New(token string) error {

	rtoken, err := redis.Get(redisPrefixRegister + u.Email)
	if err != nil {
		log.Println(err)
		return ErrTokenNotMatch
	}
	if string(rtoken) != token {
		return ErrTokenNotMatch
	}

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
	_, err = insUser.Exec(u.Username, passwordhashed, u.Email, runToken)
	if err != nil {
		return err
	}
	err = redis.Delete(redisPrefixRegister + u.Email)

	return err
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
	es := url.QueryEscape(email)
	ts := url.QueryEscape(token)
	msg := "To reset your password please click this link" +
		" https://www.yxi.io/password_new?email=" + es + "&token=" + ts

	id, err := mailgun.SimpleMessage("Reset password", msg, email)
	if err != nil {
		return
	}
	log.Println("send-mail id :", id)
	return
}

// UpdatePassByToken update password use email and token
func UpdatePassByToken(email, token, pass string) (err error) {
	rtoken, err := redis.Get(redisPrefixResetPassEmail + email)
	if err != nil {
		log.Println(err)
		return ErrTokenNotMatch
	}
	if string(rtoken) != token {
		return ErrTokenNotMatch
	}
	passwordhashed, err := crypto.HashPassword(pass)
	if err != nil {
		return err
	}
	updatePassStmt, err := mysql.Db.Prepare("UPDATE user SET password=? where email=?")
	defer updatePassStmt.Close()
	if err != nil {
		return err
	}

	_, err = updatePassStmt.Exec(passwordhashed, email)
	if err != nil {
		return err
	}
	err = redis.Delete(redisPrefixResetPassEmail + email)

	return
}
