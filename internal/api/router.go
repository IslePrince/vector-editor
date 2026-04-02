package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"vector-editor/internal/config"
	"vector-editor/internal/engine"
	"vector-editor/internal/render"
	"vector-editor/internal/storage"
)

// Server holds all dependencies for the HTTP API.
type Server struct {
	cfg    *config.Config
	store  storage.Storage
	ink    *engine.Inkscape
	worker *render.Worker
	router chi.Router
}

// NewServer wires up the API server with all dependencies.
func NewServer(cfg *config.Config, store storage.Storage, ink *engine.Inkscape, worker *render.Worker) *Server {
	s := &Server{
		cfg:    cfg,
		store:  store,
		ink:    ink,
		worker: worker,
	}
	s.router = s.buildRouter()
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) buildRouter() chi.Router {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(Logger)
	r.Use(CORS)

	r.Route("/api/v1", func(r chi.Router) {
		// Health & capabilities
		r.Get("/health", s.handleHealth)
		r.Get("/capabilities", s.handleCapabilities)

		// Documents
		r.Post("/documents", s.handleCreateDocument)
		r.Get("/documents", s.handleListDocuments)

		r.Route("/documents/{id}", func(r chi.Router) {
			r.Get("/", s.handleGetDocument)
			r.Patch("/", s.handleUpdateDocument)
			r.Delete("/", s.handleDeleteDocument)

			// Layers
			r.Post("/layers", s.handleAddLayer)
			r.Get("/layers", s.handleListLayers)
			r.Get("/layers/{lid}", s.handleGetLayer)
			r.Patch("/layers/{lid}", s.handleUpdateLayer)
			r.Delete("/layers/{lid}", s.handleDeleteLayer)

			// Renders
			r.Post("/renders", s.handleCreateRender)
			r.Get("/renders/{rid}", s.handleGetRender)
			r.Get("/renders/{rid}/download", s.handleDownloadRender)

			// SVG operations
			r.Post("/export-svg", s.handleExportSVG)
			r.Post("/text-to-path", s.handleTextToPath)
			r.Post("/boolean", s.handleBooleanOp)
		})
	})

	return r
}
