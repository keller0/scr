package handle

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/model"
	"github.com/keller0/yxi-back/util"
)

// NewCode create a new code snippet
// 200 400 500
func NewCode(c *gin.Context) {
	// get user id from jwt
	var err error
	var code model.Code
	if err = c.ShouldBindJSON(&code); err == nil {
		userid, err := util.JwtGetUserID(c.Request)
		if err != nil {
			// anonymous
			err = code.NewAnonymous()
		} else {
			code.UserID = userid
			err = code.New()
		}
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create code failed"})
		} else {
			c.String(http.StatusOK, "create code succeeded")
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusNotFound, gin.H{"error": id + " dose not exist"})
		c.Abort()
		return
	}
	var code model.Code
	code.ID = codeid

	// get user id encase the codeid's code is private
	// the error now can be ignored, because publlic code did not need auth
	code.UserID, _ = util.JwtGetUserID(c.Request)

	switch part {
	case "/content":
		content, err := code.GetCodeContentByID()
		if err != nil {
			// maybe the code need auth
			fmt.Println(part, err.Error())
			if err == model.ErrNotAllowed {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			}
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": codeRes})
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": id + " dose not exist"})
	}

}

// GetOnesCode return ones code
func GetOnesCode(c *gin.Context) {
	var err error
	var codes []model.CodeRes
	userid := c.Params.ByName("userid")
	codetype := c.DefaultQuery("type", "public")
	switch codetype {
	case "public":
		codes, err = model.GetOnesPublicCode(userid)
	case "private":
		userid, err := util.JwtGetUserID(c.Request)
		if err != nil {
			c.AbortWithError(http.StatusForbidden, err)
			return
		}
		var code model.Code
		code.UserID = userid
		codes, err = code.GetOnesPrivateCode()
	case "all":
		userid, err := util.JwtGetUserID(c.Request)
		if err != nil {
			c.AbortWithError(http.StatusForbidden, err)
			return
		}
		var code model.Code
		code.UserID = userid
		codes, err = code.GetOnesCode()
	default:
		codes, err = model.GetAllPublicCode()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"codes": codes,
	})
}

// GetCode return code list depend on param.type
func GetCode(c *gin.Context) {
	var err error
	var codes []model.CodeRes
	codetype := c.DefaultQuery("type", "public")
	switch codetype {
	case "public":
		codes, err = model.GetAllPublicCode()
	case "popular":
		codes, err = model.GetPouplarCode()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "type not supported"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	var tokenUserID int64
	if err = c.ShouldBindJSON(&code); err == nil {
		fmt.Println(code)
		tokenUserID, err = util.JwtGetUserID(c.Request)
		fmt.Println(tokenUserID)
		if err != nil {
			// anonymous user can not update code
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		err = code.UpdateCode(tokenUserID)
		if err != nil {
			fmt.Println(err)
			if err == model.ErrNotAllowed {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			}
		} else {
			c.String(http.StatusOK, "update code succeeded")
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
