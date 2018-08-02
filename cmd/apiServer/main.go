package main

import (
	"io"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/handler"
	mid "github.com/keller0/yxi-back/middleware"
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

		public.GET("/code/:codeid/*part", handle.GetCodePart)
		public.POST("/code", handle.NewCode)

		// get user's code list
		public.GET("/user/:userid/code", handle.GetOnesCode)

		public.POST("/user", handle.SendRegisterEmail)
		public.POST("/login", handle.Login)
		public.POST("/account/password/email", handle.SendResetPassEmail)
		public.POST("/account/password", handle.UpdatePassByEmail)
		public.POST("/account/complete", handle.RegisterComplete)

		run := public.Group("/run")
		{
			run.GET("/", handle.AllVersion)
			run.GET("/:language", handle.VersionsOfOne)

			rQueue := run.Group("/")
			rQueue.Use(mid.PublicLimit())
			{
				rQueue.POST("/:language", handle.RunCode)
				rQueue.POST("/:language/:version", handle.RunCode)
			}
		}
	}

	authorized := r.Group("/v1")
	authorized.Use(mid.JwtAuth())
	{
		authorized.PUT("/likes/:codeid", handle.LikeCode)
		authorized.PUT("/code", handle.UpdateCode)
		authorized.DELETE("/code/:codeid", handle.DeleteCode)
	}

	r.Run(yxiPort)

}
