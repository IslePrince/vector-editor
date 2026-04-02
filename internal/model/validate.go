package model

import (
	"fmt"
	"strings"
)

// ValidateDocument checks that a document has required fields and valid values.
func ValidateDocument(d *Document) error {
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("document name is required")
	}
	if d.Width <= 0 {
		return fmt.Errorf("width must be positive")
	}
	if d.Height <= 0 {
		return fmt.Errorf("height must be positive")
	}
	switch d.Unit {
	case UnitPx, UnitMm, UnitIn, "":
		// ok, empty defaults to px
	default:
		return fmt.Errorf("invalid unit %q, must be px, mm, or in", d.Unit)
	}
	return nil
}

// ValidateLayer checks that a layer has required fields and valid values.
func ValidateLayer(l *Layer) error {
	if strings.TrimSpace(l.Name) == "" {
		return fmt.Errorf("layer name is required")
	}
	switch l.Type {
	case LayerRect, LayerCircle, LayerEllipse, LayerPath, LayerText,
		LayerPolygon, LayerLine, LayerGroup, LayerImage:
		// ok
	case "":
		return fmt.Errorf("layer type is required")
	default:
		return fmt.Errorf("invalid layer type %q", l.Type)
	}
	if l.Opacity < 0 || l.Opacity > 1 {
		return fmt.Errorf("opacity must be between 0 and 1")
	}
	return nil
}

// ValidateRenderProfile checks that a profile name is known.
func ValidateRenderProfile(profile string) error {
	switch profile {
	case "png_hires", "png_web", "pdf_print", "svg_source":
		return nil
	default:
		return fmt.Errorf("unknown render profile %q", profile)
	}
}

// ValidateBooleanOp checks that a boolean operation name is valid.
func ValidateBooleanOp(op string) error {
	switch op {
	case "union", "intersection", "difference", "exclusion":
		return nil
	default:
		return fmt.Errorf("unknown boolean operation %q", op)
	}
}
