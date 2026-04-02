package storage

import (
	"os"
	"path/filepath"
)

// EnsureDataDirs creates the required data directory structure.
func EnsureDataDirs(baseDir string) error {
	dirs := []string{
		filepath.Join(baseDir, "documents"),
		filepath.Join(baseDir, "renders"),
		filepath.Join(baseDir, "temp"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}

// TempDir returns the path for temporary files.
func TempDir(baseDir string) string {
	return filepath.Join(baseDir, "temp")
}

// RenderDir returns the path for render output.
func RenderDir(baseDir string, docID string) string {
	return filepath.Join(baseDir, "renders", docID)
}
