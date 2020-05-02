package work

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"sync"
	"testing"
)

func TestBufferedBatch(t *testing.T) {
	var (
		m sync.Mutex

		input = "abcdefghijklmnopqrstuvwxyz"

		batch   []string
		results []string

		// Run at most 3 batches at a time
		concurrency = 3
		// Batch size of 3
		batchSize = 3

		// Counter for the batch
		batchCount int

		// Expected number of batches
		numBatches = (len(input) + batchSize - 1) / batchSize
	)

	// Build expected results
	var expectedResults = make(map[string]bool)
	for ii := 0; ii < numBatches; ii++ {
		start := ii * batchSize
		end := (ii + 1) * batchSize
		if end >= len(input) {
			end = len(input)
		}
		s := input[start:end]
		expectedResults[s] = true
	}
	require.Len(t, expectedResults, numBatches)

	b := NewBufferedBatch(concurrency, batchSize)

	for ii, ch := range input {
		s := string(ch)

		// Add an item to the batch
		batch = append(batch, s)

		// Check if the current batch is filled or if it's the last element and then start the batch job
		if b.Full() || ii == len(input)-1 {
			batchCount++

			// Define a job to run on the batch
			fn := func(ctx context.Context) error {
				processed := strings.Join(batch, "")
				fmt.Printf("Processing batch %d started: %s\n", batchCount, processed)
				randSleepMillisecond(100, 200)

				m.Lock()
				results = append(results, processed)
				m.Unlock()
				return nil
			}

			// Start the job
			err := b.Run(context.Background(), fn)
			require.NoError(t, err)

			// Remember to reset the slice for the next batch
			batch = []string{}
		}
	}

	// Wait for all jobs to finish
	errs := b.Wait()
	require.Empty(t, errs)

	fmt.Println("Results", results)

	// Compare the results
	require.Len(t, results, len(expectedResults))
	resultMap := make(map[string]bool)
	for _, r := range results {
		resultMap[r] = true
	}
	require.Equal(t, expectedResults, resultMap)

}
