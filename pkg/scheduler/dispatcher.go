package scheduler

import (
	"context"
	"sync"
)

type Dispatcher struct {
	WorkerCount int
	BuildQueue  chan Build
	workers     []*Worker
	ctx         context.Context
	cancel      context.CancelFunc
	wg          *sync.WaitGroup
}

func NewDispatcher(workerCount int, queueSize int) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &Dispatcher{
		WorkerCount: workerCount,
		BuildQueue:  make(chan Build, queueSize),
		ctx:         ctx,
		cancel:      cancel,
		wg:          &sync.WaitGroup{},
	}
}

func (d *Dispatcher) Start() {
	for i := 1; i <= d.WorkerCount; i++ {
		worker := NewWorker(i, d.BuildQueue, d.ctx)
		d.workers = append(d.workers, worker)
		worker.Start()
	}
}

func (d *Dispatcher) SubmitBuild(build Build) {
	d.wg.Add(1)
	go func() {
		d.BuildQueue <- build
		d.wg.Done()
	}()
}

func (d *Dispatcher) Stop() {
	d.cancel()
	d.wg.Wait()
	close(d.BuildQueue)
}
