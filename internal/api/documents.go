package api

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"vector-editor/internal/engine"
	"vector-editor/internal/model"
)

func (s *Server) handleCreateDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string         `json:"name"`
		Width      float64        `json:"width"`
		Height     float64        `json:"height"`
		Unit       model.Unit     `json:"unit"`
		ViewBox    *model.ViewBox `json:"viewbox,omitempty"`
		Background string         `json:"background"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid JSON: "+err.Error())
		return
	}

	now := time.Now()
	doc := &model.Document{
		ID:         uuid.New().String(),
		Name:       req.Name,
		Width:      req.Width,
		Height:     req.Height,
		Unit:       req.Unit,
		ViewBox:    req.ViewBox,
		Background: req.Background,
		Layers:     []*model.Layer{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if doc.Unit == "" {
		doc.Unit = model.UnitPx
	}
	if doc.Background == "" {
		doc.Background = "#ffffff"
	}

	if err := model.ValidateDocument(doc); err != nil {
		badRequest(w, err.Error())
		return
	}

	if err := s.store.CreateDocument(doc); err != nil {
		serverError(w, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, doc)
}

func (s *Server) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	docs, err := s.store.ListDocuments()
	if err != nil {
		serverError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, docs)
}

func (s *Server) handleGetDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := s.store.GetDocument(id)
	if err != nil {
		notFound(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

func (s *Server) handleUpdateDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := s.store.GetDocument(id)
	if err != nil {
		notFound(w, err.Error())
		return
	}

	var req struct {
		Name       *string        `json:"name,omitempty"`
		Width      *float64       `json:"width,omitempty"`
		Height     *float64       `json:"height,omitempty"`
		Unit       *model.Unit    `json:"unit,omitempty"`
		ViewBox    *model.ViewBox `json:"viewbox,omitempty"`
		Background *string        `json:"background,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid JSON: "+err.Error())
		return
	}

	if req.Name != nil {
		doc.Name = *req.Name
	}
	if req.Width != nil {
		doc.Width = *req.Width
	}
	if req.Height != nil {
		doc.Height = *req.Height
	}
	if req.Unit != nil {
		doc.Unit = *req.Unit
	}
	if req.ViewBox != nil {
		doc.ViewBox = req.ViewBox
	}
	if req.Background != nil {
		doc.Background = *req.Background
	}

	if err := model.ValidateDocument(doc); err != nil {
		badRequest(w, err.Error())
		return
	}

	if err := s.store.UpdateDocument(doc); err != nil {
		serverError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

func (s *Server) handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.store.DeleteDocument(id); err != nil {
		notFound(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleExportSVG(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := s.store.GetDocument(id)
	if err != nil {
		notFound(w, err.Error())
		return
	}

	svgBytes, err := engine.GenerateSVG(doc)
	if err != nil {
		serverError(w, "SVG generation failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+doc.Name+".svg\"")
	w.WriteHeader(http.StatusOK)
	w.Write(svgBytes)
}

func (s *Server) handleTextToPath(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := s.store.GetDocument(id)
	if err != nil {
		notFound(w, err.Error())
		return
	}

	// Generate SVG, write to temp, run Inkscape text-to-path
	svgBytes, err := engine.GenerateSVG(doc)
	if err != nil {
		serverError(w, "SVG generation failed: "+err.Error())
		return
	}

	tempDir := filepath.Join(s.cfg.DataDir, "temp")
	os.MkdirAll(tempDir, 0o755)
	inputPath := filepath.Join(tempDir, id+"_input.svg")
	outputPath := filepath.Join(tempDir, id+"_text2path.svg")

	if err := os.WriteFile(inputPath, svgBytes, 0o644); err != nil {
		serverError(w, "write temp file: "+err.Error())
		return
	}
	defer os.Remove(inputPath)

	if err := s.ink.TextToPath(r.Context(), inputPath, outputPath); err != nil {
		serverError(w, "text-to-path failed: "+err.Error())
		return
	}
	defer os.Remove(outputPath)

	result, err := os.ReadFile(outputPath)
	if err != nil {
		serverError(w, "read output: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (s *Server) handleBooleanOp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := s.store.GetDocument(id)
	if err != nil {
		notFound(w, err.Error())
		return
	}

	var req struct {
		Operation string `json:"operation"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid JSON: "+err.Error())
		return
	}
	if err := model.ValidateBooleanOp(req.Operation); err != nil {
		badRequest(w, err.Error())
		return
	}

	svgBytes, err := engine.GenerateSVG(doc)
	if err != nil {
		serverError(w, "SVG generation failed: "+err.Error())
		return
	}

	tempDir := filepath.Join(s.cfg.DataDir, "temp")
	os.MkdirAll(tempDir, 0o755)
	inputPath := filepath.Join(tempDir, id+"_bool_input.svg")
	outputPath := filepath.Join(tempDir, id+"_bool_output.svg")

	if err := os.WriteFile(inputPath, svgBytes, 0o644); err != nil {
		serverError(w, "write temp file: "+err.Error())
		return
	}
	defer os.Remove(inputPath)

	if err := s.ink.BooleanOp(r.Context(), inputPath, outputPath, req.Operation); err != nil {
		serverError(w, "boolean operation failed: "+err.Error())
		return
	}
	defer os.Remove(outputPath)

	result, err := os.ReadFile(outputPath)
	if err != nil {
		serverError(w, "read output: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
