package worker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Executor interface {
	Execute() error
	OnError(error)
	MaxRetries() int
}

type Pool struct {
	numWorkers int           // Number of workers (goroutines).
	tasks      chan Executor // Task queue (buffered channel).
	start      sync.Once     // Ensure Start() runs only once.
	stop       sync.Once     // Ensure Stop() runs only once.
	quit       chan struct{} // Signal to shutdown pool.
}

// creates a new Pool with, fixed number of workers and task queue of given size
func NewPool(numWorkers int, taskChannelSize int) *Pool {
	if numWorkers <= 0 {
		panic("numWorkers must be greater than zero")
	}
	if taskChannelSize <= 0 {
		panic("taskChannelSize must be greater than zero")
	}
	return &Pool{
		numWorkers: numWorkers,
		tasks:      make(chan Executor, taskChannelSize),
		start:      sync.Once{},
		stop:       sync.Once{},
		quit:       make(chan struct{}),
	}
}

// The intention is to run the worker goroutines here.
func (p *Pool) Start() {
	p.start.Do(func() {
		ctx := context.Background()
		p.startWorkers(ctx)
	})
}

// Gracefully stops the worker pool.
func (p *Pool) Stop() {
	p.stop.Do(func() {
		close(p.quit)
		close(p.tasks)
	})
}

// Stops the worker pool after waiting for a context timeout/cancel. Good for timeouts or graceful shutdowns in servers.
func (p *Pool) StopWithContext(ctx context.Context) {
	p.stop.Do(func() {
		close(p.quit)
		close(p.tasks)
		<-ctx.Done()
		fmt.Println("Worker pool stopped with context timeout")
	})
}

// Adds a task to the task queue. If the pool is stopped (quit is closed), it won't accept new tasks.
func (p *Pool) AddTask(t Executor) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}
func (p *Pool) startWorkers(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		go func(workerNum int) {
			// Worker Logic.
			//starting p.numWorkers number of goroutines.
			//Each goroutine is a worker with its own workerNum
			fmt.Printf("Worker number %d started\n", workerNum)
			for {
				select {
				// Handle different signals...
				//Context cancellation (ctx.Done()).
				case <-ctx.Done():
					fmt.Printf("Worker number %d stopping due to context cancellation\n", workerNum)
					return
				//Quit signal (p.quit).
				case <-p.quit:
					fmt.Printf("Worker number %d stopping due to quit signal\n", workerNum)
					return

					//New task from task queue (p.tasks).
					//The worker tries to read a task from the task queue.
					//If the channel is closed (ok == false), the worker will stop.
					//Otherwise, it will proceed to process the task.
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					//Retry Logic
					var err error
					for retry := 0; retry <= task.MaxRetries(); retry++ {
						if err = task.Execute(); err == nil {
							break
						}
						if retry <= task.MaxRetries() {
							time.Sleep(time.Duration(retry) * time.Second)
							fmt.Printf("Worker %d retrying task(attempt %d%d)\n", workerNum, retry+1, task.MaxRetries())
						}
					}
					if err != nil {
						task.OnError(err)
					}
					fmt.Printf("Worker number %d finished task\n", workerNum)
				}
			}
		}(i)
	}
}
