package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keller0/scr/cmd/apiServer/handler"
	"github.com/keller0/scr/internal/docker"
	"github.com/keller0/scr/internal/env"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	yxiPort    = env.Get("YXI_BACK_PORT", ":8090")
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
	docker.StartManagers()

	srv := &http.Server{Addr: yxiPort, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	docker.JobStop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
	os.Exit(0)
}
