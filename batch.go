package work

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Batch struct {
	wg     sync.WaitGroup
	sem    *semaphore.Weighted
	task   TaskFunc
	errors chan error
}

func NewBatch(concurrency int64) *Batch {
	return &Batch{sem: semaphore.NewWeighted(concurrency), errors: make(chan error, concurrency)}
}

func (b *Batch) Run(ctx context.Context, taskFunc TaskFunc) error {
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

		if err := fn(ctx); err != nil {
			b.errors <- err
		}
	}(taskFunc)

	return nil
}

func (b *Batch) Wait() error {
	var errMsgs []string
	b.wg.Wait()
	close(b.errors)
	for err := range b.errors {
		errMsgs = append(errMsgs, err.Error())
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf("errors detected in the multitask: [%s]", strings.Join(errMsgs, "|"))
	}
	return nil
}
