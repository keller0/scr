package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/handler"
	"github.com/keller0/yxi-back/middleware"
)

var (
	yxiPort    = os.Getenv("YXI_BACK_PORT")
	ginMode    = os.Getenv(gin.ENV_GIN_MODE)
	ginLogPath = os.Getenv("GIN_LOG_PATH")
)

func main() {

	if ginMode == gin.ReleaseMode {
		gin.DisableConsoleColor()
		f, error := os.Create(ginLogPath)
		if error != nil {
			panic("create log file failed")
		}
		gin.DefaultWriter = io.MultiWriter(f)
	}

	r := gin.Default()

	api := r.Group("/v1")
	{

		api.GET("public", handle.PublicCode)
		api.GET("public/:userid", handle.OnesPublicCode)
		api.GET("populer", handle.PopulerCode)

		api.POST("/register", handle.Register)
		api.POST("/login", handle.Login)

		p := api.Group("private").Use(mid.JwtAuth())
		{
			p.GET("/", handle.PrivateCode)
		}
	}

	r.Run(yxiPort)

}
