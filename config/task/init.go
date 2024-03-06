package task

import (
	"github.com/hibiken/asynq"
	"healthy/config/config"
)

var AsynqClient *asynq.Client
var AsynqServer *asynq.Server

func Init() asynq.RedisClientOpt {
	info := asynq.RedisClientOpt{
		Addr:     "127.0.0.1:6379",
		Password: "",
	}
	if config.Config.IsSet("redis.host") && config.Config.IsSet("redis.port") {
		info.Addr = config.Config.GetString("redis.host") + config.Config.GetString("redis.port")
	}
	if config.Config.IsSet("redis.pass") {
		info.Password = config.Config.GetString("redis.pass")
	}

	return info
}
