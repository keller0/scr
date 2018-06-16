package model

import (
	"fmt"
	"log"

	"github.com/keller0/yxi-back/db"
)

// Liked check if a user allread liked "a" code
func Liked(userID, CodeID int64) bool {
	var count int64
	err := mysql.Db.QueryRow("SELECT count(user_id) FROM likes WHERE user_id=? AND code_id=?",
		userID, CodeID).Scan(&count)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}

	return count != 0
}

// Like add a record for user like code
func Like(userID, CodeID int64) error {

	insUser, err := mysql.Db.Prepare("INSERT INTO likes(user_id, code_id) values(?,?)")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, e := insUser.Exec(userID, CodeID)
	return e
}
