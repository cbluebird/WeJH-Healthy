package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"healthy/app/midwares"
	"healthy/app/services/taskService"
	"healthy/config/config"
	"healthy/config/router"
	"healthy/config/task"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		r := gin.Default()
		r.Use(cors.Default())
		r.Use(midwares.ErrHandler())
		r.NoMethod(midwares.HandleNotFound)
		r.NoRoute(midwares.HandleNotFound)
		router.Init(r)

		if err := r.Run(":" + config.Config.GetString("port")); err != nil {
			log.Fatal("gin Server Start Failed", err)
		}
	}()

	go func() {
		mux, err := taskService.Init()
		if err != nil {
			log.Fatal(err)
		}
		if err = task.AsynqServer.Run(mux); err != nil {
			log.Fatal("mq server start failed", err)
		}
	}()

	<-stop
}
