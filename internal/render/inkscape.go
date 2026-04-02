package render

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"vector-editor/internal/engine"
	"vector-editor/internal/model"
)

// Renderer handles converting documents to output formats via Inkscape.
type Renderer struct {
	ink     *engine.Inkscape
	dataDir string
}

// NewRenderer creates a Renderer backed by the given Inkscape binary.
func NewRenderer(ink *engine.Inkscape, dataDir string) *Renderer {
	return &Renderer{ink: ink, dataDir: dataDir}
}

// RenderDocument exports a document using the given profile. It writes the SVG
// to a temp file, invokes Inkscape, and returns the output path.
func (r *Renderer) RenderDocument(ctx context.Context, doc *model.Document, profile Profile) (string, error) {
	// Generate SVG source
	svgBytes, err := engine.GenerateSVG(doc)
	if err != nil {
		return "", fmt.Errorf("generate svg: %w", err)
	}

	// Write SVG to temp input file
	renderDir := filepath.Join(r.dataDir, "renders", doc.ID)
	if err := os.MkdirAll(renderDir, 0o755); err != nil {
		return "", fmt.Errorf("create render dir: %w", err)
	}

	inputPath := filepath.Join(renderDir, "input.svg")
	if err := os.WriteFile(inputPath, svgBytes, 0o644); err != nil {
		return "", fmt.Errorf("write input svg: %w", err)
	}

	// For svg_source profile, just return the SVG directly
	if profile.ExportType == "svg" {
		outputPath := filepath.Join(renderDir, "output.svg")
		if err := os.WriteFile(outputPath, svgBytes, 0o644); err != nil {
			return "", fmt.Errorf("write output svg: %w", err)
		}
		return outputPath, nil
	}

	// Export via Inkscape
	ext := profile.ExportType
	outputPath := filepath.Join(renderDir, "output."+ext)
	opts := ProfileToExportOptions(profile, inputPath, outputPath)

	if err := r.ink.Export(ctx, opts); err != nil {
		return "", fmt.Errorf("inkscape export: %w", err)
	}

	return outputPath, nil
}
