package handle

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

func PrivateCode(c *gin.Context) {

	c.String(http.StatusOK, "ssssssssssssssssssss")

}
