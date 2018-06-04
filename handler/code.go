package handle

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/middleware"
)

func PrivateCode(c *gin.Context) {
	token := c.GetHeader("Authorization")

	id, err := mid.JwtGetUserID(token)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
	}

	c.String(http.StatusOK, fmt.Sprintf("user's id : %d", id))

}
