package files_test

import (
	"github.com/devinbarry/crev/internal/files"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

// TestGetContentMapOfFiles tests reading the content of files and handling empty directories,
// verifying that directory contents are properly handled and empty directories are marked.
func TestGetContentMapOfFiles(t *testing.T) {
	rootDir := t.TempDir()

	// Create directory structure
	subDir1 := filepath.Join(rootDir, "subdir_1")
	subDir2 := filepath.Join(rootDir, "subdir_2")
	createDir(t, subDir1)
	createDir(t, subDir2)

	// Create files
	createFile(t, filepath.Join(rootDir, "file1.txt"), "content1")
	createFile(t, filepath.Join(subDir1, "file2.txt"), "content2")

	filePaths := []string{
		filepath.Join(rootDir, "file1.txt"),
		subDir1,
		filepath.Join(subDir1, "file2.txt"),
		subDir2,
	}

	fileContentMap, err := files.GetContentMapOfFiles(filePaths, 10)
	require.NoError(t, err, "GetContentMapOfFiles failed")

	// Verify file contents
	require.Equal(t, "content1", fileContentMap[filepath.Join(rootDir, "file1.txt")], "Incorrect content for file1.txt")
	require.Equal(t, "content2", fileContentMap[filepath.Join(subDir1, "file2.txt")], "Incorrect content for file2.txt")

	// Verify directory handling
	require.NotContains(t, fileContentMap, subDir1, "Non-empty directory should not be in map")
	require.Equal(t, "empty directory", fileContentMap[subDir2], "Empty directory not properly marked")
}
