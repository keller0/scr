package main

import (
	"io"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/handler"
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
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	public := r.Group("/v1")
	{
		// get code list
		public.GET("/code", handle.GetCode)

		public.GET("/code/content/:codeid", handle.GetCodeContent)
		public.POST("/code", handle.NewCode)

		public.GET("/user/:userid/code", handle.GetOnesCode)

		public.POST("/register", handle.Register)
		public.POST("/login", handle.Login)

	}

	r.Run(yxiPort)

}
