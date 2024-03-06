package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"healthy/app/midwares"
	"healthy/app/services/taskService"
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

		err := r.Run()
		if err != nil {
			log.Fatal("ServerStartFailed", err)
		}
	}()

	go func() {
		info := task.Init()
		task.AsynqClient = asynq.NewClient(info)
		task.AsynqServer = asynq.NewServer(
			info,
			asynq.Config{Concurrency: 10}, //Concurrency表示最大并发处理任务数。
		)
		if err := task.AsynqServer.Run(asynq.HandlerFunc(taskService.Handler)); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
}
