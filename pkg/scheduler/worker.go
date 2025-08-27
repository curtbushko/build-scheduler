package scheduler

import (
	"context"
	"fmt"
)

type Worker struct {
	ID         int
	buildQueue <-chan Build
	context    context.Context
}

func NewWorker(id int, buildChan <-chan Build, ctx context.Context) *Worker {
	return &Worker{
		ID:         id,
		buildQueue: buildChan,
		context:    ctx,
	}
}

func (w *Worker) Start() {
	go func() {
		fmt.Println("Worker started", w.ID)
		for {
			select {
			case <-w.context.Done():
				fmt.Println("Stopping worker", w.ID)
				return
			case build := <-w.buildQueue:
				build.Process()
			}
		}
	}()
}
