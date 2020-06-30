package workerpool

import (
	"sync"
	"sync/atomic"
)

type job struct {
	id uint64
	f  func()
}

// WorkerPool executes tasks from a queue. The queue size and the concurrency
// can be configured. Use NewWorkerPool() to create a WorkerPool instance.
type WorkerPool struct {
	jobsQueueCh chan (*job)
	wg          sync.WaitGroup
	idSequence  uint64
}

func (p *WorkerPool) work() {
	for j := range p.jobsQueueCh {
		j.f()
		p.wg.Done()
	}
}

// NewWorkerPool constructs a new WorkerPool. The concurrency parameter
// describes how many concurrent goroutines will be executing the tasks. The
// queueSize parameter indicates how many tasks can be inserted via AddTask()
// before the insertion starts to block.
func NewWorkerPool(concurrency, queueSize int) *WorkerPool {
	p := &WorkerPool{
		jobsQueueCh: make(chan (*job), queueSize),
		idSequence:  0,
	}

	for i := 0; i < concurrency; i++ {
		go p.work()
	}

	return p
}

// AddTask inserts a new task into the queue. This method may block if
// the queue is full. The method returns the job id assigned to the task.
func (p *WorkerPool) AddTask(f func()) uint64 {
	p.wg.Add(1)

	id := atomic.AddUint64(&p.idSequence, 1)
	p.jobsQueueCh <- &job{id: id, f: f}

	return id
}

// WaitAll waits until all the tasks put into the queue have finished.
func (p *WorkerPool) WaitAll() {
	p.wg.Wait()
}

// Close shuts down the pool immediately. Tasks left in the queue will be
// discarded.
func (p *WorkerPool) Close() {
	close(p.jobsQueueCh)
}
