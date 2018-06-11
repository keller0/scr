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

// GetCodeContent return code with content
func GetCodeContent(c *gin.Context) {
	var err error
	var content string
	codeid := c.Params.ByName("codeid")
	var code model.Code
	code.ID, err = strconv.ParseInt(codeid, 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": codeid + " dose not exist"})
		c.Abort()
		return
	}

	// get user id encase the codeid's code is private
	code.UserID, _ = util.JwtGetUserID(c.Request)

	content, err = code.GetCodeContent()
	if err != nil {
		if err == model.ErrNotAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": model.ErrNotAllowed.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": codeid + " dose not exist"})
		}
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"content": content,
	})
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
