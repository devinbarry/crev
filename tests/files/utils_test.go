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

// assertFileSetMatches verifies that the actual file paths exactly match expectations.
// Parameters:
//   - t: the testing context
//   - actualPaths: the slice of file paths returned by GetAllFilePaths
//   - expectedPaths: the slice of file paths that should be included
//   - notExpectedPaths: the slice of file paths that should be excluded
//   - msgAndArgs: optional message and arguments for test failure output
func assertFileSetMatches(t *testing.T, actualPaths []string, expectedPaths []string, notExpectedPaths []string, msgAndArgs ...interface{}) {
	t.Helper()

	// First check that we have exactly the expected number of files
	require.Len(t, actualPaths, len(expectedPaths),
		"Wrong number of files returned. Expected exactly %d files, got %d. %v",
		len(expectedPaths), len(actualPaths), msgAndArgs)

	// Check that all expected paths are present
	require.ElementsMatch(t, expectedPaths, actualPaths,
		"File paths don't match expected set. %v", msgAndArgs)

	// Check that none of the excluded paths are present
	for _, excludedPath := range notExpectedPaths {
		require.NotContains(t, actualPaths, excludedPath,
			"Found excluded path %q in results. %v", excludedPath, msgAndArgs)
	}
}
