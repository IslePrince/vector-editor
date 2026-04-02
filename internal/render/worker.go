package render

import (
	"context"
	"log"
	"sync"
	"time"

	"vector-editor/internal/model"
)

// JobFunc is a function that resolves a document by ID for rendering.
type JobFunc func(documentID string) (*model.Document, error)

// Worker processes render jobs from a channel.
type Worker struct {
	renderer *Renderer
	jobs     chan *model.RenderJob
	mu       sync.RWMutex
	results  map[string]*model.RenderJob
	docFunc  JobFunc
}

// NewWorker creates a render worker pool.
func NewWorker(renderer *Renderer, maxWorkers int, docFunc JobFunc) *Worker {
	w := &Worker{
		renderer: renderer,
		jobs:     make(chan *model.RenderJob, 100),
		results:  make(map[string]*model.RenderJob),
		docFunc:  docFunc,
	}
	for i := 0; i < maxWorkers; i++ {
		go w.process()
	}
	return w
}

// Submit enqueues a render job for processing.
func (w *Worker) Submit(job *model.RenderJob) {
	w.mu.Lock()
	w.results[job.ID] = job
	w.mu.Unlock()
	w.jobs <- job
}

// GetJob returns the current state of a render job.
func (w *Worker) GetJob(jobID string) (*model.RenderJob, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	j, ok := w.results[jobID]
	return j, ok
}

func (w *Worker) process() {
	for job := range w.jobs {
		w.mu.Lock()
		job.Status = model.RenderProcessing
		job.Progress = 10
		job.UpdatedAt = time.Now()
		w.mu.Unlock()

		doc, err := w.docFunc(job.DocumentID)
		if err != nil {
			w.failJob(job, "document not found: "+err.Error())
			continue
		}

		w.mu.Lock()
		job.Progress = 30
		w.mu.Unlock()

		profile, ok := Profiles[job.Profile]
		if !ok {
			w.failJob(job, "unknown profile: "+job.Profile)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		outputPath, err := w.renderer.RenderDocument(ctx, doc, profile)
		cancel()

		if err != nil {
			w.failJob(job, err.Error())
			continue
		}

		w.mu.Lock()
		job.Status = model.RenderComplete
		job.Progress = 100
		job.OutputPath = outputPath
		job.UpdatedAt = time.Now()
		w.mu.Unlock()

		log.Printf("render complete: job=%s doc=%s profile=%s", job.ID, job.DocumentID, job.Profile)
	}
}

func (w *Worker) failJob(job *model.RenderJob, errMsg string) {
	w.mu.Lock()
	job.Status = model.RenderFailed
	job.Error = errMsg
	job.UpdatedAt = time.Now()
	w.mu.Unlock()
	log.Printf("render failed: job=%s error=%s", job.ID, errMsg)
}
