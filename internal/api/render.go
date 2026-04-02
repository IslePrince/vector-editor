package api

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"vector-editor/internal/model"
)

func (s *Server) handleCreateRender(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")
	if _, err := s.store.GetDocument(docID); err != nil {
		notFound(w, err.Error())
		return
	}

	var req struct {
		Profile string `json:"profile"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid JSON: "+err.Error())
		return
	}
	if err := model.ValidateRenderProfile(req.Profile); err != nil {
		badRequest(w, err.Error())
		return
	}

	now := time.Now()
	job := &model.RenderJob{
		ID:         uuid.New().String(),
		DocumentID: docID,
		Profile:    req.Profile,
		Status:     model.RenderPending,
		Progress:   0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	s.worker.Submit(job)
	writeJSON(w, http.StatusAccepted, job)
}

func (s *Server) handleGetRender(w http.ResponseWriter, r *http.Request) {
	renderID := chi.URLParam(r, "rid")
	job, ok := s.worker.GetJob(renderID)
	if !ok {
		notFound(w, "render job not found")
		return
	}
	writeJSON(w, http.StatusOK, job)
}

func (s *Server) handleDownloadRender(w http.ResponseWriter, r *http.Request) {
	renderID := chi.URLParam(r, "rid")
	job, ok := s.worker.GetJob(renderID)
	if !ok {
		notFound(w, "render job not found")
		return
	}

	if job.Status != model.RenderComplete {
		badRequest(w, "render is not yet complete, status: "+string(job.Status))
		return
	}

	if job.OutputPath == "" {
		serverError(w, "output path is empty")
		return
	}

	data, err := os.ReadFile(job.OutputPath)
	if err != nil {
		serverError(w, "read output file: "+err.Error())
		return
	}

	ext := filepath.Ext(job.OutputPath)
	contentType := "application/octet-stream"
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".pdf":
		contentType = "application/pdf"
	case ".svg":
		contentType = "image/svg+xml"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\"render"+ext+"\"")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
