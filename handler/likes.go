package handle

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi.io/model"
)

// LikeCode create
// 200 401 404 409 500
func LikeCode(c *gin.Context) {
	codeID := c.Param("codeid")
	codeid, err := strconv.ParseInt(codeID, 10, 64)
	if err != nil {
		fmt.Println(codeid, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["Bad Requset"]})
		return
	}
	// check if logined
	uid, e := c.Get("uid")
	if !e {
		c.JSON(http.StatusUnauthorized, gin.H{"errNumber": responseErr["Like Code Not Allowed"]})
		return
	}
	userid := uid.(int64)
	// check code exist
	if !model.CodeExist(codeid) {
		c.JSON(http.StatusNotFound, gin.H{"errNumber": responseErr["CodeNotExist"]})
		return
	}
	// check if already liked
	if model.Liked(userid, codeid) {
		c.JSON(http.StatusConflict, gin.H{"errNumber": responseErr["Already Liked"]})
		return
	}
	err = model.Like(userid, codeid)
	if err != nil {
		fmt.Println(userid, codeid, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["ServerErr Like Code Failed"]})
		return
	}
	c.String(http.StatusOK, "like code succeeded")
}
