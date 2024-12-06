package files_test

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a file with specified content.
// Ensures that the file is created successfully.
func createFile(t *testing.T, path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err, "Failed to create file %s", path)
}

// Helper function to create a directory.
// Ensures that the directory is created successfully.
func createDir(t *testing.T, path string) {
	err := os.MkdirAll(path, 0755)
	require.NoError(t, err, "Failed to create directory %s", path)
}

// Helper function to create multiple files in a directory structure.
// Takes a map of relative paths to content and creates all files under rootDir.
func createFiles(t *testing.T, rootDir string, files map[string]string) {
	for path, content := range files {
		fullPath := filepath.Join(rootDir, path)
		dir := filepath.Dir(fullPath)
		createDir(t, dir)
		createFile(t, fullPath, content)
	}
}
