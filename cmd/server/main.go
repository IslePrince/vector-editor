package main

import (
	"fmt"
	"log"
	"net/http"

	"vector-editor/internal/api"
	"vector-editor/internal/config"
	"vector-editor/internal/engine"
	"vector-editor/internal/model"
	"vector-editor/internal/render"
	"vector-editor/internal/storage"
)

func main() {
	cfg := config.Load()

	// Ensure data directories exist
	if err := storage.EnsureDataDirs(cfg.DataDir); err != nil {
		log.Fatalf("failed to create data dirs: %v", err)
	}

	// Storage
	store, err := storage.NewLocalStorage(cfg.DataDir)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	// Inkscape engine
	ink := engine.NewInkscape(cfg.InkscapePath)

	// Renderer + worker pool
	renderer := render.NewRenderer(ink, cfg.DataDir)
	docFunc := func(docID string) (*model.Document, error) {
		return store.GetDocument(docID)
	}
	worker := render.NewWorker(renderer, cfg.MaxWorkers, docFunc)

	// HTTP server
	srv := api.NewServer(cfg, store, ink, worker)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("vector-editor API starting on %s", addr)
	log.Printf("data directory: %s", cfg.DataDir)

	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
