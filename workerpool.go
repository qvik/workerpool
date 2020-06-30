package workerpool

import (
	"sync"
	"sync/atomic"
)

type task struct {
	id uint64
	f  func() error
}

// Result is the result of a task; if Error is set, the task returned an error.
type Result struct {
	TaskID uint64
	Error  error
}

// WorkerPool executes tasks from a queue. The queue size and the concurrency
// can be configured. Use NewWorkerPool() to create a WorkerPool instance.
type WorkerPool struct {
	taskQueueCh chan (*task)
	resultCh    chan (*Result)
	wg          sync.WaitGroup
	idSequence  uint64
}

// GetResultsChannel returns the Results channel. Consume this channel
// until it is closed to get all the results for your tasks.
func (p *WorkerPool) GetResultsChannel() chan (*Result) {
	return p.resultCh
}

func (p *WorkerPool) work() {
	for t := range p.taskQueueCh {
		err := t.f()

		if p.resultCh != nil {
			// Emit the result via the channel; this might block if results
			// channel is full.
			p.resultCh <- &Result{TaskID: t.id, Error: err}
		}

		p.wg.Done()
	}
}

// NewWorkerPool constructs a new WorkerPool. The concurrency parameter
// describes how many concurrent goroutines will be executing the tasks. The
// queueSize parameter indicates how many tasks can be inserted via AddTask()
// before the insertion starts to block. The resultsSize parameter
// indicates the size of the Results channel (see GetResultsChannel()). This
// value should be set to the number of tasks going to be run; however, it
// can be set to a lower value as well, but then the channel MUST be consumed
// before calling WaitAll() or the pool will get stuck. Specify 0 to ignore
// the results (Results wont get emitted on the channel).
func NewWorkerPool(concurrency, queueSize, resultsSize int) *WorkerPool {
	p := &WorkerPool{
		taskQueueCh: make(chan (*task), queueSize),
		idSequence:  0,
	}

	if resultsSize > 0 {
		p.resultCh = make(chan (*Result), resultsSize)
	}

	for i := 0; i < concurrency; i++ {
		go p.work()
	}

	return p
}

// AddTask inserts a new task into the queue. This method may block if
// the queue is full. The method returns the job id assigned to the task.
func (p *WorkerPool) AddTask(f func() error) uint64 {
	p.wg.Add(1)

	id := atomic.AddUint64(&p.idSequence, 1)
	p.taskQueueCh <- &task{id: id, f: f}

	return id
}

// WaitAll waits until all the tasks put into the queue have finished.
// Another (finer grained) option for tracking the completion of your tasks is
// to read the Results channel returned by GetResultsChannel().
func (p *WorkerPool) WaitAll() {
	p.wg.Wait()
}

// Close shuts down the pool immediately. Tasks left in the queue will be
// discarded. Note that the results from the Results channel
// (see GetResultsChannel()) can still be read.
func (p *WorkerPool) Close() {
	close(p.taskQueueCh)
}
