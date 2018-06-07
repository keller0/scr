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

	public := r.Group("/v1")
	{

		public.GET("/code/pub", handle.PublicCode)
		public.GET("/code/top", handle.PopulerCode)
		public.GET("/code/pub/:userid", handle.OnesPublicCode)
		public.POST("/code/new", handle.NewCode)

		public.POST("/register", handle.Register)
		public.POST("/login", handle.Login)

	}

	private := r.Group("/v1").Use(mid.JwtAuth())
	{
		// get one's private  code
		private.GET("/code/private", handle.PrivateCode)
	}

	r.Run(yxiPort)

}
