package work

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type lockedCount struct {
	sync.Mutex
	count int64
}

func (c *lockedCount) Add(n int64) {
	c.Lock()
	c.count += n
	c.Unlock()
}

type Results struct {
	combinedNames string
	calls         int
}

// Test that multiple tasks run successfully
func TestBatch(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	names := []string{
		"calvin",
		"pk",
		"anthony",
		"abdul",
		"nirav",
		"bing",
		"ryan",
		"anson",
		"matt",
		"tony",
	}
	var concurrency int64 = 3
	results := &Results{}
	batch := NewBatch(concurrency)

	// currentlyRunning is a counter tracking the number of currently running tasks
	currentlyRunning := lockedCount{}
	// onStart is used to signal when a new task starts
	onStart := make(chan struct{}, len(names))

	lock := sync.Mutex{}
	for _, name := range names {
		name := name
		err := batch.Run(ctx, func(ctx context.Context) error {
			onStart <- struct{}{}
			currentlyRunning.Add(1)
			defer currentlyRunning.Add(-1)
			randSleepMillisecond(1, 10)

			lock.Lock()
			results.combinedNames += "|" + name
			results.calls += 1
			lock.Unlock()
			println(results.combinedNames)
			if name == names[len(names)-1] {
				close(onStart)
			}
			return nil
		})
		require.NoError(t, err)
	}

	// Check whenever a routine starts that there are only at most `concurrency` number of tasks running
	go func() {
		<-onStart
		require.LessOrEqual(t, currentlyRunning.count, concurrency)
	}()

	// Wait for any currently running tasks to finish processing
	err := batch.Wait()

	require.NoError(t, err)
	require.Equal(t, len(names), results.calls)
	for _, name := range names {
		require.Contains(t, results.combinedNames, name)
	}
}

func randSleepMillisecond(min, max int) {
	time.Sleep(time.Millisecond*time.Duration(rand.Intn(max)) + time.Millisecond*time.Duration(min))
}
