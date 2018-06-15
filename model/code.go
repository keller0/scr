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

// CodeRes include code info, code's user, and likes
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

// GetOnesCode get one user's code
func (c *Code) GetOnesCode() ([]CodeRes, error) {

	userid := strconv.Itoa(int(c.UserID))
	return getCodes("code.user_id="+userid, "code.create_at", "desc", "15")
}

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

// GetPouplarCode get all populer code from code table
func GetPouplarCode() ([]CodeRes, error) {
	return getCodes("code.public=true", "likes", "desc", "15")
}

// this does not contain code's content
func getCodes(where, orderby, order, limit string) ([]CodeRes, error) {
	selOut, err := mysql.Db.Query(
		"SELECT code.id, IFNULL(user.username,\"" + anonymousUser + "\") username, " +
			"code.title, code.description, code.lang, code.filename," +
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
		var username, title, lang, description, filename, createtat, updateat string
		var pub bool
		err := selOut.Scan(&id, &username, &title, &description, &lang,
			&filename, &createtat, &updateat, &pub, &likes)
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
		code.CreateAt = createtat
		code.UpdateAt = updateat
		code.Public = pub
		codes = append(codes, code)
	}
	return codes, nil
}

// GetCodeByID use code's id get info
func (c *Code) GetCodeByID() (CodeRes, error) {
	var code CodeRes
	var userid int64
	err := mysql.Db.QueryRow(
		"SELECT code.id, IFNULL(user.username,\""+anonymousUser+"\") username,"+
			"IFNULL(code.user_id, 0), code.content,"+
			"code.create_at, code.update_at,"+
			"code.title, code.description, code.lang, code.filename,"+
			"code.public, count(likes.code_id) likes "+
			"FROM code left join user on code.user_id=user.id "+
			"left join likes on likes.code_id = code.id "+
			"where code.id=?", c.ID).Scan(&code.ID, &code.UserName,
		&userid, &code.Content, &code.CreateAt, &code.UpdateAt, &code.Title,
		&code.Description, &code.Lang, &code.FileName, &code.Public, &code.Likes)

	if err != nil {
		return CodeRes{}, err
	}

	// if code is not public, the userid need matche
	if !code.Public && userid != c.UserID {
		return CodeRes{}, ErrNotAllowed
	}
	return code, nil
}

// GetCodeContentByID use code's id get  content
func (c *Code) GetCodeContentByID() (string, error) {
	var content string
	var isPublic bool
	var userid int64
	err := mysql.Db.QueryRow(
		"SELECT content, IFNULL(user_id, 0), public FROM code WHERE id=?",
		c.ID).Scan(&content, &userid, &isPublic)
	// code does not exist
	if err != nil {
		return "", err
	}
	// if code is not public, the userid need matche
	if !isPublic && userid != c.UserID {
		return "", ErrNotAllowed
	}

	return content, nil
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
	_, err = insCode.Exec(c.Title, c.Description, c.Lang, c.FileName, c.Content, true)
	return err
}
