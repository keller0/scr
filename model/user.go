package model

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/keller0/yxi-back/db"
	"github.com/keller0/yxi-back/util"
)

var (
	reUsername          = regexp.MustCompile("^[a-zA-Z0-9-_.]+$")
	errUserLoginInvalid = errors.New("Invalid User Login")
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

// Validate validates the required fields and formats.
func (u *User) Validate() error {
	switch {
	case len(u.Username) == 0:
		return errUserLoginInvalid
	case len(u.Username) > 250:
		return errUserLoginInvalid
	case !reUsername.MatchString(u.Username):
		return errUserLoginInvalid
	default:
		return nil
	}
}

// UsernameExist check if username already existed
func (u *User) UsernameExist() bool {

	var id int64
	err := mysql.Db.QueryRow("SELECT id FROM user WHERE username=?", u.Username).Scan(&id)
	if err != nil {
		return true
	}
	return id != 0
}

// EmailExist check if email already existed
func (u *User) EmailExist() bool {

	var id int64
	err := mysql.Db.QueryRow("SELECT id FROM user WHERE email=?", u.Email).Scan(&id)
	if err != nil {
		return true
	}
	return id != 0
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

	if util.CheckPasswordHash(u.Password, password) {
		exp := time.Now().Add(time.Hour * 1).Unix()
		tokenString, err := util.JwtGenToken(id, u.Username, runToken, exp)
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}

	return "", ErrLoginWrongPass
}

// New create a new user account
func (u *User) New() error {

	var runToken = util.RandStringRunes(40)
	passwordhashed, err := util.HashPassword(u.Password)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	insUser, err := mysql.Db.Prepare("INSERT INTO user(username, password, email, run_token) values(?,?,?,?)")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, e := insUser.Exec(u.Username, passwordhashed, u.Email, runToken)
	return e
}
