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

	p := NewWorkerPool(concurrency, 100)

	for i := 0; i < numTasks; i++ {
		p.AddTask(func() {
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
		})
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
