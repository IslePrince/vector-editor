package queue

import (
	"fmt"
	"sync"

	"vector-editor/internal/model"
)

// MemoryQueue is an in-memory job queue.
type MemoryQueue struct {
	mu   sync.RWMutex
	jobs map[string]*model.RenderJob
	ch   chan *model.RenderJob
}

// NewMemoryQueue creates a new in-memory queue.
func NewMemoryQueue(bufferSize int) *MemoryQueue {
	return &MemoryQueue{
		jobs: make(map[string]*model.RenderJob),
		ch:   make(chan *model.RenderJob, bufferSize),
	}
}

// Enqueue adds a job to the queue.
func (q *MemoryQueue) Enqueue(job *model.RenderJob) error {
	q.mu.Lock()
	q.jobs[job.ID] = job
	q.mu.Unlock()

	select {
	case q.ch <- job:
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

// Dequeue blocks until a job is available.
func (q *MemoryQueue) Dequeue() (*model.RenderJob, error) {
	job, ok := <-q.ch
	if !ok {
		return nil, fmt.Errorf("queue closed")
	}
	return job, nil
}

// Get returns a job by ID.
func (q *MemoryQueue) Get(id string) (*model.RenderJob, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	j, ok := q.jobs[id]
	return j, ok
}

// List returns all jobs.
func (q *MemoryQueue) List() []*model.RenderJob {
	q.mu.RLock()
	defer q.mu.RUnlock()
	result := make([]*model.RenderJob, 0, len(q.jobs))
	for _, j := range q.jobs {
		result = append(result, j)
	}
	return result
}
