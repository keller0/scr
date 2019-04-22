package main

import (
	"io"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keller0/scr/cmd/apiServer/handler"
	"github.com/keller0/scr/internal/env"
)

var (
	yxiPort    = env.Get("YXI_BACK_PORT", ":8090")
	ginMode    = env.Get("GIN_MODE", "debug")
	ginLogPath = env.Get("GIN_LOG_PATH", "/var/log/yxi/api.log")
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

	v1 := r.Group("/v1")
	{

		v1.GET("/", handler.AllVersion)
		v1.GET("/:language", handler.VersionsOfOne)

		v1.POST("/:language", handler.RunCode)
		v1.POST("/:language/:version", handler.RunCode)

	}

	r.Run(yxiPort)

}
