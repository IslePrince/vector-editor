package engine

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Inkscape wraps the Inkscape CLI binary for vector operations.
type Inkscape struct {
	BinPath string
}

// NewInkscape creates a new Inkscape engine wrapper.
func NewInkscape(binPath string) *Inkscape {
	return &Inkscape{BinPath: binPath}
}

// ExportOptions configures an export operation.
type ExportOptions struct {
	InputPath  string
	OutputPath string
	ExportType string  // png, pdf, svg, eps
	DPI        int
	Width      int
	Height     int
}

// Export converts an SVG file to the specified format using Inkscape CLI.
func (ink *Inkscape) Export(ctx context.Context, opts ExportOptions) error {
	args := []string{
		opts.InputPath,
		"--export-type=" + opts.ExportType,
		"--export-filename=" + opts.OutputPath,
	}
	if opts.DPI > 0 {
		args = append(args, fmt.Sprintf("--export-dpi=%d", opts.DPI))
	}
	if opts.Width > 0 {
		args = append(args, fmt.Sprintf("--export-width=%d", opts.Width))
	}
	if opts.Height > 0 {
		args = append(args, fmt.Sprintf("--export-height=%d", opts.Height))
	}

	cmd := exec.CommandContext(ctx, ink.BinPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("inkscape export failed: %w\noutput: %s", err, string(output))
	}
	return nil
}

// TextToPath converts all text objects in an SVG to paths using Inkscape.
func (ink *Inkscape) TextToPath(ctx context.Context, inputPath, outputPath string) error {
	actions := "select-all;object-to-path;export-filename:" + outputPath + ";export-do"
	args := []string{
		inputPath,
		"--actions=" + actions,
	}

	cmd := exec.CommandContext(ctx, ink.BinPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("inkscape text-to-path failed: %w\noutput: %s", err, string(output))
	}
	return nil
}

// BooleanOp performs a boolean path operation on two selected objects.
// op must be one of: union, intersection, difference, exclusion.
func (ink *Inkscape) BooleanOp(ctx context.Context, inputPath, outputPath, op string) error {
	var actionName string
	switch op {
	case "union":
		actionName = "path-union"
	case "intersection":
		actionName = "path-intersection"
	case "difference":
		actionName = "path-difference"
	case "exclusion":
		actionName = "path-exclusion"
	default:
		return fmt.Errorf("unsupported boolean operation: %s", op)
	}

	actions := strings.Join([]string{
		"select-all",
		actionName,
		"export-filename:" + outputPath,
		"export-do",
	}, ";")

	args := []string{
		inputPath,
		"--actions=" + actions,
	}

	cmd := exec.CommandContext(ctx, ink.BinPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("inkscape boolean %s failed: %w\noutput: %s", op, err, string(output))
	}
	return nil
}

// Version returns the installed Inkscape version string.
func (ink *Inkscape) Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, ink.BinPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("inkscape version check failed: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
