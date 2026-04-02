package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"vector-editor/internal/model"
)

func (s *Server) handleAddLayer(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")
	if _, err := s.store.GetDocument(docID); err != nil {
		notFound(w, err.Error())
		return
	}

	var req struct {
		Name      string            `json:"name"`
		Type      model.LayerType   `json:"type"`
		Visible   *bool             `json:"visible,omitempty"`
		Locked    *bool             `json:"locked,omitempty"`
		Opacity   *float64          `json:"opacity,omitempty"`
		Fill      *model.Fill       `json:"fill,omitempty"`
		Stroke    *model.Stroke     `json:"stroke,omitempty"`
		Transform *model.Transform  `json:"transform,omitempty"`
		Data      map[string]string `json:"data,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid JSON: "+err.Error())
		return
	}

	now := time.Now()
	layer := &model.Layer{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Type:      req.Type,
		Visible:   true,
		Locked:    false,
		Opacity:   1.0,
		Fill:      req.Fill,
		Stroke:    req.Stroke,
		Transform: req.Transform,
		Data:      req.Data,
		Children:  []*model.Layer{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.Visible != nil {
		layer.Visible = *req.Visible
	}
	if req.Locked != nil {
		layer.Locked = *req.Locked
	}
	if req.Opacity != nil {
		layer.Opacity = *req.Opacity
	}

	if err := model.ValidateLayer(layer); err != nil {
		badRequest(w, err.Error())
		return
	}

	if err := s.store.AddLayer(docID, layer); err != nil {
		serverError(w, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, layer)
}

func (s *Server) handleListLayers(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")
	layers, err := s.store.ListLayers(docID)
	if err != nil {
		notFound(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, layers)
}

func (s *Server) handleGetLayer(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")
	layerID := chi.URLParam(r, "lid")

	layer, err := s.store.GetLayer(docID, layerID)
	if err != nil {
		notFound(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, layer)
}

func (s *Server) handleUpdateLayer(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")
	layerID := chi.URLParam(r, "lid")

	layer, err := s.store.GetLayer(docID, layerID)
	if err != nil {
		notFound(w, err.Error())
		return
	}

	var req struct {
		Name      *string           `json:"name,omitempty"`
		Type      *model.LayerType  `json:"type,omitempty"`
		Visible   *bool             `json:"visible,omitempty"`
		Locked    *bool             `json:"locked,omitempty"`
		Opacity   *float64          `json:"opacity,omitempty"`
		Fill      *model.Fill       `json:"fill,omitempty"`
		Stroke    *model.Stroke     `json:"stroke,omitempty"`
		Transform *model.Transform  `json:"transform,omitempty"`
		Data      map[string]string `json:"data,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid JSON: "+err.Error())
		return
	}

	if req.Name != nil {
		layer.Name = *req.Name
	}
	if req.Type != nil {
		layer.Type = *req.Type
	}
	if req.Visible != nil {
		layer.Visible = *req.Visible
	}
	if req.Locked != nil {
		layer.Locked = *req.Locked
	}
	if req.Opacity != nil {
		layer.Opacity = *req.Opacity
	}
	if req.Fill != nil {
		layer.Fill = req.Fill
	}
	if req.Stroke != nil {
		layer.Stroke = req.Stroke
	}
	if req.Transform != nil {
		layer.Transform = req.Transform
	}
	if req.Data != nil {
		if layer.Data == nil {
			layer.Data = make(map[string]string)
		}
		for k, v := range req.Data {
			layer.Data[k] = v
		}
	}
	layer.UpdatedAt = time.Now()

	if err := s.store.UpdateLayer(docID, layer); err != nil {
		serverError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, layer)
}

func (s *Server) handleDeleteLayer(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")
	layerID := chi.URLParam(r, "lid")

	if err := s.store.DeleteLayer(docID, layerID); err != nil {
		notFound(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
