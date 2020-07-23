package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keller0/scr/cmd/apiServer/handler"
	"github.com/keller0/scr/internal/docker"
	"github.com/keller0/scr/internal/env"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	yxiPort    = env.Get("YXI_BACK_PORT", ":8090")
	yxiHost    = env.Get("YXI_BACK_HOST", "localhost")
	ginLogPath = env.Get("GIN_LOG_PATH", "api.log")
)

func main() {
	log.Info("starting...")
	configLog()

	r := configEngine()
	v1 := r.Group("/v1")
	{
		v1.GET("/runners", handler.AllRunners)
		v1.GET("/runners/:language", handler.VersionsOfOne)

		v1.POST("/:language", handler.RunCode)
		v1.POST("/:language/:version", handler.RunCode)
	}

	docker.StartManagers()

	srv := &http.Server{Addr: yxiHost + yxiPort, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Info("listening at ", yxiHost+yxiPort)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server...")

	docker.JobStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown with error:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
	os.Exit(0)
}

func configLog() {

	logFile, err := os.OpenFile(ginLogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatal("open log file failed:", err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

}

func configEngine() *gin.Engine {
	r := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))
	r.Use(gin.Recovery())
	return r
}
