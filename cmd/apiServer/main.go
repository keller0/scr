package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keller0/scr/cmd/apiServer/handler"
	"github.com/keller0/scr/internal/env"
	log "github.com/sirupsen/logrus"
)

var (
	yxiPort = env.Get("YXI_BACK_PORT", ":8090")
	ginLogPath = env.Get("GIN_LOG_PATH", "/var/log/yxi/api.log")
)

func main() {
	log.Info("starting...")
	r := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))
	r.Use(gin.Recovery())
	v1 := r.Group("/v1")
	{

		v1.GET("/", handler.AllVersion)
		v1.GET("/:language", handler.VersionsOfOne)

		v1.POST("/:language", handler.RunCode)
		v1.POST("/:language/:version", handler.RunCode)

	}

	err := r.Run(yxiPort)
	log.Error(err)

}
