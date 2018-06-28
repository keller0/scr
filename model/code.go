package model

import (
	"fmt"
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
func (c *Code) GetOnesCode(offset string) ([]CodeRes, error) {

	userid := strconv.Itoa(int(c.UserID))
	return getCodes("code.user_id="+userid, "code.create_at", "desc", "15", offset)
}

// GetOnesPublicCode get one user's public code
func GetOnesPublicCode(userid, offset string) ([]CodeRes, error) {
	return getCodes("code.public=true and code.user_id="+userid, "code.create_at", "desc", "15", offset)
}

// GetOnesPrivateCode get one user's private code
func (c *Code) GetOnesPrivateCode(offset string) ([]CodeRes, error) {

	userid := strconv.Itoa(int(c.UserID))
	return getCodes("code.public=false and code.user_id="+userid, "code.create_at", "desc", "15", offset)
}

// GetAllPublicCode get all public code from code table
func GetAllPublicCode(offset string) ([]CodeRes, error) {
	return getCodes("code.public=true", "code.create_at", "desc", "15", offset)
}

// GetPouplarCode get all populer code from code table
func GetPouplarCode(offset string) ([]CodeRes, error) {
	return getCodes("code.public=true", "likes", "desc", "15", offset)
}

// this does not contain code's content
func getCodes(where, orderby, order, limit, offset string) ([]CodeRes, error) {
	fmt.Println(where, orderby, order, limit, offset)
	selOut, err := mysql.Db.Query(
		"SELECT code.id, IFNULL(user.username,\"" + anonymousUser + "\") username, " +
			"code.title, code.description, code.lang, code.filename," +
			"code.create_at, code.update_at, code.public, count(likes.code_id) likes " +
			"FROM code left join user on code.user_id=user.id " +
			"left join likes on likes.code_id = code.id " +
			"where " + where + " group by code.id ORDER BY " + orderby + " " + order +
			" LIMIT " + limit + " OFFSET " + offset)
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

// CodeExist check if code exist use code id
func CodeExist(id int64) bool {
	var count int64
	err := mysql.Db.QueryRow("SELECT count(id) FROM code WHERE id=? ",
		id).Scan(&count)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}

	return count != 0
}

// UpdateCode if anonymous the userid is 0, save as null
func (c *Code) UpdateCode(tokenUserID int64) error {

	var theuser int64
	var err error
	err = mysql.Db.QueryRow("SELECT IFNULL(user_id, -1) FROM code WHERE id=? ",
		c.ID).Scan(&theuser)
	if err != nil {
		return err
	}
	if theuser != tokenUserID {
		return ErrNotAllowed
	}

	updateCode, err := mysql.Db.Prepare("UPDATE code " +
		"SET title=?, description=?, filename=?, content=?, public=?, user_id=?, update_at=now() WHERE id=?")

	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	if c.UserID == 0 {
		_, err = updateCode.Exec(c.Title, c.Description, c.FileName, c.Content, c.Public, nil, c.ID)
	} else if c.UserID == tokenUserID {
		_, err = updateCode.Exec(c.Title, c.Description, c.FileName, c.Content, c.Public, c.UserID, c.ID)
	} else {
		return ErrNotAllowed
	}

	return err
}
