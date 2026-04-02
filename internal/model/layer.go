package model

import "time"

// LayerType identifies the kind of vector element.
type LayerType string

const (
	LayerRect    LayerType = "rect"
	LayerCircle  LayerType = "circle"
	LayerEllipse LayerType = "ellipse"
	LayerPath    LayerType = "path"
	LayerText    LayerType = "text"
	LayerPolygon LayerType = "polygon"
	LayerLine    LayerType = "line"
	LayerGroup   LayerType = "group"
	LayerImage   LayerType = "image"
)

// FillType identifies the kind of fill.
type FillType string

const (
	FillSolid    FillType = "solid"
	FillGradient FillType = "gradient"
	FillNone     FillType = "none"
)

// Fill describes the interior coloring of a shape.
type Fill struct {
	Type    FillType `json:"type"`
	Color   string   `json:"color"`
	Opacity float64  `json:"opacity"`
}

// Stroke describes the outline of a shape.
type Stroke struct {
	Color     string  `json:"color"`
	Width     float64 `json:"width"`
	LineCap   string  `json:"line_cap"`
	LineJoin  string  `json:"line_join"`
	DashArray string  `json:"dash_array"`
	Opacity   float64 `json:"opacity"`
}

// Transform describes an SVG transform.
type Transform struct {
	Translate *[2]float64 `json:"translate,omitempty"`
	Rotate    *float64    `json:"rotate,omitempty"`
	Scale     *[2]float64 `json:"scale,omitempty"`
	SkewX     *float64    `json:"skew_x,omitempty"`
	SkewY     *float64    `json:"skew_y,omitempty"`
	Matrix    *[6]float64 `json:"matrix,omitempty"`
}

// Layer represents a single vector element or group in the document.
type Layer struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      LayerType         `json:"type"`
	Visible   bool              `json:"visible"`
	Locked    bool              `json:"locked"`
	Opacity   float64           `json:"opacity"`
	Fill      *Fill             `json:"fill,omitempty"`
	Stroke    *Stroke           `json:"stroke,omitempty"`
	Transform *Transform        `json:"transform,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
	Children  []*Layer          `json:"children,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
