package render

import "vector-editor/internal/engine"

// Profile defines settings for a render output format.
type Profile struct {
	Name       string
	ExportType string
	DPI        int
}

// Profiles maps profile names to their settings.
var Profiles = map[string]Profile{
	"png_hires": {Name: "png_hires", ExportType: "png", DPI: 300},
	"png_web":   {Name: "png_web", ExportType: "png", DPI: 72},
	"pdf_print": {Name: "pdf_print", ExportType: "pdf", DPI: 300},
	"svg_source": {Name: "svg_source", ExportType: "svg", DPI: 0},
}

// ProfileToExportOptions converts a render profile to Inkscape export options.
func ProfileToExportOptions(p Profile, inputPath, outputPath string) engine.ExportOptions {
	return engine.ExportOptions{
		InputPath:  inputPath,
		OutputPath: outputPath,
		ExportType: p.ExportType,
		DPI:        p.DPI,
	}
}

// ProfileNames returns all available profile names.
func ProfileNames() []string {
	return []string{"png_hires", "png_web", "pdf_print", "svg_source"}
}
