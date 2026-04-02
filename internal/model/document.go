package model

import "time"

// Unit represents an SVG coordinate unit.
type Unit string

const (
	UnitPx Unit = "px"
	UnitMm Unit = "mm"
	UnitIn Unit = "in"
)

// ViewBox describes the SVG viewBox attribute.
type ViewBox struct {
	MinX   float64 `json:"min_x"`
	MinY   float64 `json:"min_y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Document represents a single SVG document.
type Document struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Width      float64   `json:"width"`
	Height     float64   `json:"height"`
	Unit       Unit      `json:"unit"`
	ViewBox    *ViewBox  `json:"viewbox,omitempty"`
	Background string    `json:"background"`
	Layers     []*Layer  `json:"layers"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Capabilities describes what this editor supports.
type Capabilities struct {
	Shapes         []string `json:"shapes"`
	RenderProfiles []string `json:"render_profiles"`
	BooleanOps     []string `json:"boolean_ops"`
	Units          []string `json:"units"`
}
