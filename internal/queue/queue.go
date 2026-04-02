package queue

import "vector-editor/internal/model"

// Queue defines the interface for job queuing.
type Queue interface {
	Enqueue(job *model.RenderJob) error
	Dequeue() (*model.RenderJob, error)
	Get(id string) (*model.RenderJob, bool)
	List() []*model.RenderJob
}
