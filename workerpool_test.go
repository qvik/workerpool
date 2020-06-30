package workerpool

import (
	"sync"
	"testing"
	"time"
)

func TestConcurrency(t *testing.T) {
	concurrency := 10
	numTasks := concurrency * 2
	perceivedConcurrency := 0
	maxConcurrency := 0
	invocations := 0

	var rwlock sync.RWMutex

	p := NewWorkerPool(concurrency, 100, numTasks)

	for i := 0; i < numTasks; i++ {
		p.AddTask(func() error {
			rwlock.Lock()
			invocations++
			perceivedConcurrency++
			if perceivedConcurrency > maxConcurrency {
				maxConcurrency = perceivedConcurrency
			}
			rwlock.Unlock()

			// Process our "task"
			time.Sleep(250 * time.Millisecond)

			rwlock.Lock()
			perceivedConcurrency--
			rwlock.Unlock()

			return nil
		})
	}

	// Read all results
	for i := 0; i < numTasks; i++ {
		res := <-p.GetResultsChannel()

		if res.Error != nil {
			t.Errorf("task failed unexpectedly")
		}
	}

	p.WaitAll()
	p.Close()

	if invocations != numTasks {
		t.Errorf("incorrect number of invocations: %v", invocations)
	}

	if maxConcurrency != concurrency {
		t.Errorf("concurrency differed from specified: %v", maxConcurrency)
	}
}
