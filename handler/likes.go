package handle

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/model"
	"github.com/keller0/yxi-back/util"
)

// LikeCode create
// 200 401 404 409 500
func LikeCode(c *gin.Context) {
	codeID := c.Param("codeid")
	codeid, err := strconv.ParseInt(codeID, 10, 64)
	if err != nil {
		fmt.Println(codeid, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": codeID + " dose not exist"})
		c.Abort()
		return
	}
	// check if logined
	userid, err := util.JwtGetUserID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "please login"})
		c.Abort()
		return
	}
	// check code exist
	if !model.CodeExist(codeid) {
		c.JSON(http.StatusNotFound, gin.H{"error": codeID + " code dose not exist"})
		c.Abort()
		return
	}
	// check if already liked
	if model.Liked(userid, codeid) {
		c.JSON(http.StatusConflict, gin.H{"error": "already liked"})
		c.Abort()
		return
	}
	err = model.Like(userid, codeid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.String(http.StatusOK, "like code succeeded")
}
