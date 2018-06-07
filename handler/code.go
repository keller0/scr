package handle

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/model"
	"github.com/keller0/yxi-back/util"
)

// PrivateCode return one's private code
func PrivateCode(c *gin.Context) {

	// get user id from jwt
	userid, err := util.JwtGetUserID(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
	}
	var code model.Code
	code.UserID = userid
	codes, err := code.GetOnesPrivateCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"codes": codes})

}

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

// PublicCode return all public code
func PublicCode(c *gin.Context) {
	codes, err := model.GetAllPublicCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"codes": codes,
	})

}

// OnesPublicCode return one's public code
func OnesPublicCode(c *gin.Context) {
	userid := c.Params.ByName("userid")
	codes, err := model.GetOnesPublicCode(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"codes": codes,
	})
}

// PopulerCode return most liked code
func PopulerCode(c *gin.Context) {

	codes, err := model.GetPouplerCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"codes": codes,
	})
}
