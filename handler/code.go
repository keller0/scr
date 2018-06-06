package handle

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/model"
	"github.com/keller0/yxi-back/util"
)

func PrivateCode(c *gin.Context) {
	token := c.GetHeader("Authorization")

	id, err := util.JwtGetUserID(token)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
	}

	c.String(http.StatusOK, fmt.Sprintf("user's id : %d", id))

}
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
