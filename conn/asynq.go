package conn

import (
	"ems/config"
	"github.com/hibiken/asynq"
)

var asynqClient *asynq.Client
var asynqInspector *asynq.Inspector

func InitAsynqClient() {
	asynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     config.Asynq().RedisAddr,
		Password: config.Asynq().Password,
		DB:       config.Asynq().DB,
	})
}

func InitAsynqInspector() {
	asynqInspector = asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     config.Asynq().RedisAddr,
		Password: config.Asynq().Password,
		DB:       config.Asynq().DB,
	})
}

func Asynq() *asynq.Client {
	return asynqClient
}

func AsynqInspector() *asynq.Inspector {
	return asynqInspector
}
