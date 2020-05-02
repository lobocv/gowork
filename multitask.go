package work

import (
	"context"
	"sync"
)

type TaskFunc func(ctx context.Context) error

type task struct {
	fn  TaskFunc
	ctx context.Context
}

type MultiTask struct {
	// errors is a channel to signal fatal errors while processing legs of the job
	errors chan error
	// wait is the waitgroup used to synchronize the independently running legs of the job
	wg sync.WaitGroup
	// List of tasks to run concurrently
	tasks []task
}

func (m *MultiTask) AddTask(ctx context.Context, fn TaskFunc) {
	m.tasks = append(m.tasks, task{ctx: ctx, fn: fn})
}

func (m *MultiTask) Run() []error {
	m.errors = make(chan error, len(m.tasks))
	m.wg.Add(len(m.tasks))
	for _, tt := range m.tasks {

		go func(t task) {
			defer m.wg.Done()
			if err := t.fn(t.ctx); err != nil {
				m.errors <- err
			}
		}(tt)
	}
	return m.wait()
}

// wait waits for all the asynchronous legs of the job to complete and then consolidates any errors produced
// in any leg into one overall error for the job
func (m *MultiTask) wait() (errs []error) {
	m.wg.Wait()
	close(m.errors)
	for err := range m.errors {
		errs = append(errs, err)
	}
	return
}
