package handle

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/internal/token"
	"github.com/keller0/yxi-back/model"
)

// NewCode create a new code snippet
// 200 400 500
func NewCode(c *gin.Context) {
	// get user id from jwt
	var err error
	var code model.Code
	if err = c.ShouldBindJSON(&code); err == nil {
		userid, err := token.JwtGetUserID(c.Request)
		if err != nil {
			// anonymous
			err = code.NewAnonymous()
		} else {
			code.UserID = userid
			err = code.New()
		}
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["ServerErr Create Code Failed"]})
		} else {
			c.String(http.StatusOK, "create code succeeded")
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Bad Requset"]})
	}
}

// GetCodePart return part of code
// 200 403 404
func GetCodePart(c *gin.Context) {
	id := c.Param("codeid")
	part := c.Param("part")
	codeid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println(codeid, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
		c.Abort()
		return
	}
	var code model.Code
	code.ID = codeid

	// get user id encase the codeid's code is private
	// the error now can be ignored, because publlic code did not need auth
	code.UserID, _ = token.JwtGetUserID(c.Request)

	switch part {
	case "/content":
		content, err := code.GetCodeContentByID()
		if err != nil {
			// maybe the code need auth
			fmt.Println(part, err.Error())
			if err == model.ErrNotAllowed {
				c.JSON(http.StatusForbidden, gin.H{"errNumber": responseErr["Get Code Not Allowed"]})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
			}
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"content": content})
	case "/":
		codeRes, err := code.GetCodeByID()
		if err != nil {
			fmt.Println(part, err.Error())
			if err == model.ErrNotAllowed {
				c.JSON(http.StatusForbidden, gin.H{"errNumber": responseErr["Get Code Not Allowed"]})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
			}
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": codeRes})
	default:
		c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
	}

}

// GetOnesCode return ones code
// 200 403 500
func GetOnesCode(c *gin.Context) {
	var err error
	var codes []model.CodeRes
	userid := c.Params.ByName("userid")
	codetype := c.DefaultQuery("type", "public")
	offsite := c.DefaultQuery("off", "0")
	switch codetype {
	case "public":
		codes, err = model.GetOnesPublicCode(userid, offsite)
	case "private":
		userid, err := token.JwtGetUserID(c.Request)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"errNumber": responseErr["Get Code Not Allowed"]})
			c.Abort()
			return
		}
		var code model.Code
		code.UserID = userid
		codes, err = code.GetOnesPrivateCode(offsite)
	case "all":
		userid, err := token.JwtGetUserID(c.Request)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"errNumber": responseErr["Get Code Not Allowed"]})
			c.Abort()
			return
		}
		var code model.Code
		code.UserID = userid
		codes, err = code.GetOnesCode(offsite)
	default:
		codes, err = model.GetAllPublicCode(offsite)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["ServerErr Get Code Failed"]})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"codes": codes,
	})
}

// GetCode return code list depend on param.type
// 200 404 500
func GetCode(c *gin.Context) {
	var err error
	var codes []model.CodeRes
	codetype := c.DefaultQuery("type", "public")
	offsite := c.DefaultQuery("off", "0")
	switch codetype {
	case "public":
		codes, err = model.GetAllPublicCode(offsite)
	case "popular":
		codes, err = model.GetPouplarCode(offsite)
	default:
		c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["ServerErr Get Code Failed"]})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"codes": codes,
	})

}

// UpdateCode update code
// 200 400 401 403 404
func UpdateCode(c *gin.Context) {

	var err error
	var code model.Code

	if err = c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Bad Requset"]})
		c.Abort()
		return
	}
	tokenUserID, e := c.Get("uid")
	if !e {
		// anonymous user can not update code
		c.JSON(http.StatusUnauthorized, gin.H{"errNumber": responseErr["Update Code Not Allowed"]})
		c.Abort()
		return
	}
	err = code.UpdateCode(tokenUserID.(int64))
	if err != nil {
		fmt.Println(err)
		if err == model.ErrNotAllowed {
			c.JSON(http.StatusForbidden, gin.H{"errNumber": responseErr["Update Code Not Allowed"]})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
		}
	} else {
		c.String(http.StatusOK, "update code succeeded")
	}

}

// DeleteCode delete code by id
func DeleteCode(c *gin.Context) {
	id := c.Param("codeid")

	codeid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println(codeid, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
		c.Abort()
		return
	}
	tokenUserID, e := c.Get("uid")
	if !e {
		c.JSON(http.StatusUnauthorized, gin.H{"errNumber": responseErr["Delete Code Not Allowed"]})
		c.Abort()
		return
	}

	err = model.DeleteCodeByID(codeid, tokenUserID.(int64))
	if err != nil {
		if err == model.ErrUserNotMatch {
			c.JSON(http.StatusUnauthorized, gin.H{"errNumber": responseErr["Delete Code Not Allowed"]})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
		}
	} else {
		c.String(http.StatusOK, "delete code succeeded")
	}
}
