package util

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Job interface {
	Run()
}

type jobWrapper struct {
	action func()
}

func (j jobWrapper) Run() {
	j.action()
}

type WorkerPool struct {
	name   string
	worker int
	jobs   chan Job
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewWorkerPool(name string, worker int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		name:   name,
		worker: worker,
		jobs:   make(chan Job, worker),
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

func (w *WorkerPool) Run() {
	w.wg.Add(w.worker)
	for i := 0; i < w.worker; i++ {
		go func() {
			defer w.wg.Done()
			for {
				select {
				case j := <-w.jobs:
					j.Run()
				case <-w.ctx.Done():
					return
				}
			}
		}()
	}
}

func (w *WorkerPool) Push(action func()) {
	w.jobs <- jobWrapper{
		action: action,
	}
}

func (w *WorkerPool) Shutdown(d time.Duration) error {
	shutdown := make(chan struct{})
	w.cancel()
	go func() {
		w.wg.Wait()
		shutdown <- struct{}{}
	}()
	select {
	case <-time.After(d):
		return errors.New(fmt.Sprintf("shutdown worker pool %s with time out", w.name))
	case <-shutdown:
		return nil
	}
}
