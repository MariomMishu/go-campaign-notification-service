package worker

import (
	"ems/config"
	"fmt"
	"github.com/hibiken/asynq"
)

func StartAsynqWorker(mux *asynq.ServeMux) {
	worker := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.Asynq().RedisAddr,
			DB:       config.Asynq().DB,
			Password: config.Asynq().Password,
		},
		asynq.Config{
			Concurrency: config.Asynq().Concurrency,
			Queues: map[string]int{
				config.Asynq().Queue: 1,
			},
		},
	)
	if err := worker.Run(mux); err != nil {
		panic(fmt.Sprint("Could not Start worker: %v", err))
	}
}
