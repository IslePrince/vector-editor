package api

import (
	"net/http"

	"vector-editor/internal/model"
	"vector-editor/internal/render"
)

// HealthResponse is the health check payload.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "vector-editor",
		Version: "1.0.0",
	})
}

func (s *Server) handleCapabilities(w http.ResponseWriter, r *http.Request) {
	caps := model.Capabilities{
		Shapes:         []string{"rect", "circle", "ellipse", "path", "text", "polygon", "line", "group", "image"},
		RenderProfiles: render.ProfileNames(),
		BooleanOps:     []string{"union", "intersection", "difference", "exclusion"},
		Units:          []string{"px", "mm", "in"},
	}
	writeJSON(w, http.StatusOK, caps)
}
