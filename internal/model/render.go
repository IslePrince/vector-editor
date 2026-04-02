package model

import "time"

// RenderStatus represents the state of a render job.
type RenderStatus string

const (
	RenderPending    RenderStatus = "pending"
	RenderProcessing RenderStatus = "processing"
	RenderComplete   RenderStatus = "complete"
	RenderFailed     RenderStatus = "failed"
)

// RenderJob tracks an asynchronous render/export operation.
type RenderJob struct {
	ID         string       `json:"id"`
	DocumentID string       `json:"document_id"`
	Profile    string       `json:"profile"`
	Status     RenderStatus `json:"status"`
	Progress   int          `json:"progress"`
	OutputPath string       `json:"output_path,omitempty"`
	Error      string       `json:"error,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}
