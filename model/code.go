package model

import (
	"log"
	"strconv"

	"github.com/keller0/yxi-back/db"
)

// Code for every code snippet
type Code struct {
	ID int64 `json:"id"`

	UserID int64 `json:"userid"`

	Title string `json:"title"`

	Description string `json:"description"`

	Lang string `json:"lang"`

	FileName string `json:"filename"`

	Content string `json:"content"`

	Public bool `json:"public"`

	CreateAt string `json:"createat"`

	UpdateAt string `json:"updateat"`
}

type CodeRes struct {
	ID int64 `json:"id"`

	UserName string `json:"username"`

	Likes int64 `json:"likes"`

	Title string `json:"title"`

	Description string `json:"description"`

	Lang string `json:"lang"`

	FileName string `json:"filename"`

	Content string `json:"content"`

	Public bool `json:"public"`

	CreateAt string `json:"createat"`

	UpdateAt string `json:"updateat"`
}

var anonymousUser = "anonymous"

// GetOnesPublicCode get one user's public code
func GetOnesPublicCode(userid string) ([]CodeRes, error) {
	return getCodes("code.public=true and code.user_id="+userid, "code.create_at", "desc", "15")
}

// GetOnesPrivateCode get one user's private code
func (c *Code) GetOnesPrivateCode() ([]CodeRes, error) {

	userid := strconv.Itoa(int(c.UserID))
	return getCodes("code.public=false and code.user_id="+userid, "code.create_at", "desc", "15")
}

// GetAllPublicCode get all public code from code table
func GetAllPublicCode() ([]CodeRes, error) {
	return getCodes("code.public=true", "code.create_at", "desc", "15")
}

// GetPouplerCode get all populer code from code table
func GetPouplerCode() ([]CodeRes, error) {
	return getCodes("code.public=true", "create_at", "desc", "15")
}

func getCodes(where, orderby, order, limit string) ([]CodeRes, error) {
	selOut, err := mysql.Db.Query(
		"SELECT code.id, IFNULL(user.username,\"" + anonymousUser + "\") username, " +
			"code.title, code.description, code.lang, code.filename, code.content, " +
			"code.create_at, code.update_at, code.public, count(likes.code_id) likes " +
			"FROM code left join user on code.user_id=user.id " +
			"left join likes on likes.code_id = code.id " +
			"where " + where + " group by code.id ORDER BY " + orderby + " " + order + " LIMIT " + limit)
	if err != nil {
		return nil, err
	}

	codes := []CodeRes{}
	for selOut.Next() {
		code := CodeRes{}
		var id, likes int64
		var username, title, lang, description, filename, content, createtat, updateat string
		var pub bool
		err := selOut.Scan(&id, &username, &title, &description, &lang,
			&filename, &content, &createtat, &updateat, &pub, &likes)
		if err != nil {
			return nil, err
		}
		code.ID = id
		code.UserName = username
		code.Likes = likes
		code.Title = title
		code.Lang = lang
		code.Description = description
		code.FileName = filename
		code.Content = content
		code.CreateAt = createtat
		code.UpdateAt = updateat
		code.Public = pub
		codes = append(codes, code)
	}
	return codes, nil
}

// New create a new code snippet recoard
func (c *Code) New() error {
	var err error
	insCode, err := mysql.Db.Prepare("INSERT INTO code" +
		"(user_id, title, description, lang, filename, content, public) " +
		"values(?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	_, err = insCode.Exec(c.UserID, c.Title, c.Description, c.Lang, c.FileName, c.Content, c.Public)

	return err
}

// NewAnonymous create a new code snippet recoard without user_id
func (c *Code) NewAnonymous() error {
	var err error
	insCode, err := mysql.Db.Prepare("INSERT INTO code" +
		"(title, description, lang, filename, content, public) " +
		"values(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	// anonymous code could only be public
	_, err = insCode.Exec(c.Title, c.Description, c.Lang, c.FileName, c.Content, "true")
	return err
}
