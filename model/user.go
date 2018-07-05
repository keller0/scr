package model

import (
	"log"
	"time"

	"github.com/keller0/yxi-back/db/mysql"
	"github.com/keller0/yxi-back/util"
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

	var count int64
	err := mysql.Db.QueryRow("SELECT count(id) FROM user WHERE email=?", u.Email).Scan(&count)
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

	if util.CheckPasswordHash(u.Password, password) {
		exp := time.Now().Add(time.Hour * 24 * 15).Unix()
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
