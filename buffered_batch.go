package work

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type BufferedBatch struct {
	wg   sync.WaitGroup
	sem  *semaphore.Weighted
	task TaskFunc

	count   int
	bufSize int

	errors      chan error
	startSignal chan bool
}

func NewBufferedBatch(concurrency int, bufSize int) *BufferedBatch {

	b := &BufferedBatch{
		sem:         semaphore.NewWeighted(int64(concurrency)),
		bufSize:     bufSize,
		errors:      make(chan error, concurrency),
		startSignal: make(chan bool, concurrency),
	}
	return b
}

func (b *BufferedBatch) Run(ctx context.Context, taskFunc TaskFunc) error {
	err := b.sem.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	b.wg.Add(1)

	go func(fn TaskFunc) {

		defer func() {
			b.wg.Done()
			b.sem.Release(1)
		}()

		b.startSignal <- true
		if err := fn(ctx); err != nil {
			b.errors <- err
		}

	}(taskFunc)

	<-b.startSignal
	return nil
}

func (b *BufferedBatch) Full() bool {
	b.count++
	if b.count == b.bufSize {
		b.count = 0
		return true
	} else {
		return false
	}
}

func (b *BufferedBatch) Wait() (errs []error) {
	b.wg.Wait()
	close(b.errors)
	for err := range b.errors {
		errs = append(errs, err)
	}
	return
}
