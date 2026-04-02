package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"vector-editor/internal/model"
)

// LocalStorage stores documents as JSON files on disk with an in-memory index.
type LocalStorage struct {
	baseDir string
	mu      sync.RWMutex
	docs    map[string]*model.Document
}

// NewLocalStorage creates a LocalStorage rooted at baseDir and loads existing documents.
func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	docsDir := filepath.Join(baseDir, "documents")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		return nil, fmt.Errorf("create documents dir: %w", err)
	}

	s := &LocalStorage{
		baseDir: baseDir,
		docs:    make(map[string]*model.Document),
	}

	// Load existing documents
	entries, err := os.ReadDir(docsDir)
	if err != nil {
		return nil, fmt.Errorf("read documents dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(docsDir, e.Name()))
		if err != nil {
			continue
		}
		var doc model.Document
		if err := json.Unmarshal(data, &doc); err != nil {
			continue
		}
		s.docs[doc.ID] = &doc
	}

	return s, nil
}

func (s *LocalStorage) docPath(id string) string {
	return filepath.Join(s.baseDir, "documents", id+".json")
}

func (s *LocalStorage) persist(doc *model.Document) error {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal document: %w", err)
	}
	return os.WriteFile(s.docPath(doc.ID), data, 0o644)
}

// CreateDocument stores a new document.
func (s *LocalStorage) CreateDocument(doc *model.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.docs[doc.ID]; exists {
		return fmt.Errorf("document %s already exists", doc.ID)
	}
	s.docs[doc.ID] = doc
	return s.persist(doc)
}

// GetDocument retrieves a document by ID.
func (s *LocalStorage) GetDocument(id string) (*model.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, ok := s.docs[id]
	if !ok {
		return nil, fmt.Errorf("document %s not found", id)
	}
	return doc, nil
}

// ListDocuments returns all documents.
func (s *LocalStorage) ListDocuments() ([]*model.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*model.Document, 0, len(s.docs))
	for _, doc := range s.docs {
		result = append(result, doc)
	}
	return result, nil
}

// UpdateDocument replaces an existing document.
func (s *LocalStorage) UpdateDocument(doc *model.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.docs[doc.ID]; !exists {
		return fmt.Errorf("document %s not found", doc.ID)
	}
	doc.UpdatedAt = time.Now()
	s.docs[doc.ID] = doc
	return s.persist(doc)
}

// DeleteDocument removes a document.
func (s *LocalStorage) DeleteDocument(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.docs[id]; !exists {
		return fmt.Errorf("document %s not found", id)
	}
	delete(s.docs, id)
	return os.Remove(s.docPath(id))
}

// AddLayer appends a layer to a document.
func (s *LocalStorage) AddLayer(docID string, layer *model.Layer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.docs[docID]
	if !ok {
		return fmt.Errorf("document %s not found", docID)
	}
	doc.Layers = append(doc.Layers, layer)
	doc.UpdatedAt = time.Now()
	return s.persist(doc)
}

// GetLayer finds a layer by ID within a document.
func (s *LocalStorage) GetLayer(docID, layerID string) (*model.Layer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, ok := s.docs[docID]
	if !ok {
		return nil, fmt.Errorf("document %s not found", docID)
	}
	layer := findLayer(doc.Layers, layerID)
	if layer == nil {
		return nil, fmt.Errorf("layer %s not found in document %s", layerID, docID)
	}
	return layer, nil
}

// ListLayers returns all layers in a document.
func (s *LocalStorage) ListLayers(docID string) ([]*model.Layer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, ok := s.docs[docID]
	if !ok {
		return nil, fmt.Errorf("document %s not found", docID)
	}
	return doc.Layers, nil
}

// UpdateLayer replaces a layer in a document.
func (s *LocalStorage) UpdateLayer(docID string, layer *model.Layer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.docs[docID]
	if !ok {
		return fmt.Errorf("document %s not found", docID)
	}
	for i, l := range doc.Layers {
		if l.ID == layer.ID {
			doc.Layers[i] = layer
			doc.UpdatedAt = time.Now()
			return s.persist(doc)
		}
	}
	return fmt.Errorf("layer %s not found in document %s", layer.ID, docID)
}

// DeleteLayer removes a layer from a document.
func (s *LocalStorage) DeleteLayer(docID, layerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.docs[docID]
	if !ok {
		return fmt.Errorf("document %s not found", docID)
	}
	for i, l := range doc.Layers {
		if l.ID == layerID {
			doc.Layers = append(doc.Layers[:i], doc.Layers[i+1:]...)
			doc.UpdatedAt = time.Now()
			return s.persist(doc)
		}
	}
	return fmt.Errorf("layer %s not found in document %s", layerID, docID)
}

func findLayer(layers []*model.Layer, id string) *model.Layer {
	for _, l := range layers {
		if l.ID == id {
			return l
		}
		if found := findLayer(l.Children, id); found != nil {
			return found
		}
	}
	return nil
}
