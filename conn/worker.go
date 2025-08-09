package conn

import (
	"ems/config"
	"ems/worker"
)

var workerPool *worker.Pool

func ConnectWorker() {
	workerPool = worker.NewPool(config.App().NumberOfWorkers, 2*config.App().NumberOfWorkers)
	workerPool.Start()
}

func WorkerPool() *worker.Pool {
	return workerPool
}
