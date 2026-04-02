package storage

import "vector-editor/internal/model"

// Storage defines the interface for document persistence.
type Storage interface {
	// Documents
	CreateDocument(doc *model.Document) error
	GetDocument(id string) (*model.Document, error)
	ListDocuments() ([]*model.Document, error)
	UpdateDocument(doc *model.Document) error
	DeleteDocument(id string) error

	// Layers
	AddLayer(docID string, layer *model.Layer) error
	GetLayer(docID, layerID string) (*model.Layer, error)
	ListLayers(docID string) ([]*model.Layer, error)
	UpdateLayer(docID string, layer *model.Layer) error
	DeleteLayer(docID, layerID string) error
}
