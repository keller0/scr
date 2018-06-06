package model

import (
	"github.com/keller0/yxi-back/db"
)

// Code for every code
type Code struct {
	ID int64 `json:"id"`

	UserID int64 `json:"userid"`

	Title string `json:"title"`

	Description string `json:"description"`

	Lang string `json:"lang"`

	CreateTime string `json:"createtime"`

	ModifyTime string `json:"modifytime"`

	FileName string `json:"filename"`

	Content string `json:"content"`

	Public bool `json:"public"`
}

func GetOnesPublicCode(userid string) ([]Code, error) {
	return getCodes("public=true and user_id="+userid, "id", "desc", "15")
}

func GetAllPublicCode() ([]Code, error) {
	return getCodes("public=true", "id", "desc", "15")
}

func GetPouplerCode() ([]Code, error) {
	return getCodes("public=true", "create_time", "desc", "15")
}
func getCodes(where, orderby, order, limit string) ([]Code, error) {
	selOut, err := mysql.Db.Query("SELECT id, user_id,title, " +
		"description,lang,filename,content,create_time,modify_time " +
		"FROM code where " + where + " ORDER BY " + orderby + " " + order + " LIMIT " + limit)
	if err != nil {
		return nil, err
	}
	codes := []Code{}
	for selOut.Next() {
		code := Code{}
		var id, userid int64
		var title, lang, description, filename, content, createtime, modifytime string

		err := selOut.Scan(&id, &userid, &title, &description, &lang, &filename, &content, &createtime, &modifytime)
		if err != nil {
			return nil, err
		}
		code.ID = id
		code.UserID = userid
		code.Title = title
		code.Lang = lang
		code.Description = description
		code.FileName = filename
		code.Content = content
		code.CreateTime = createtime
		code.ModifyTime = modifytime
		codes = append(codes, code)
	}
	return codes, nil
}
