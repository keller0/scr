package main

import (
	"io"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi.io/handler"
	"github.com/keller0/yxi.io/internal"
)

var (
	yxiPort    = internal.GetEnv("YXI_BACK_PORT", ":8090")
	ginMode    = internal.GetEnv("GIN_MODE", "debug")
	ginLogPath = internal.GetEnv("GIN_LOG_PATH", "/var/log/yxi/api.log")
)

func main() {

	if ginMode == gin.ReleaseMode {
		gin.DisableConsoleColor()
		f, err := os.Create(ginLogPath)
		if err != nil {
			panic("create log file failed")
		}
		gin.DefaultWriter = io.MultiWriter(f)
	}

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	run := r.Group("/api")
	{
		v1 := run.Group("/v1")
		{

			v1.GET("/", handle.AllVersion)
			v1.GET("/:language", handle.VersionsOfOne)

			v1.POST("/:language", handle.RunCode)
			v1.POST("/:language/:version", handle.RunCode)

		}
	}

	r.Run(yxiPort)

}
