package files_test

import (
	"github.com/devinbarry/crev/internal/files"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

// TestGetAllFilePathsExcludeDirTrailingSlash tests that directories are correctly excluded
// regardless of whether the exclude pattern has a trailing slash or not.
func TestGetAllFilePathsExcludeDirTrailingSlash(t *testing.T) {
	rootDir := t.TempDir()

	// Create directories and files
	dirPath := filepath.Join(rootDir, "dir")
	err := os.Mkdir(dirPath, 0755)
	require.NoError(t, err, "Failed to create directory")

	fileInDir := filepath.Join(dirPath, "file.txt")
	err = os.WriteFile(fileInDir, []byte("content"), 0644)
	require.NoError(t, err, "Failed to create file")

	// Test excluding directory without trailing slash
	excludePatterns := []string{"dir"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed without trailing slash")
	require.Empty(t, filePaths, "Expected no files when excluding directory without slash")

	// Test excluding directory with trailing slash
	excludePatterns = []string{"dir/"}
	filePaths, err = files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed with trailing slash")
	require.Empty(t, filePaths, "Expected no files when excluding directory with slash")
}

// TestGetAllFilePathsExcludeFileVsDirectory tests that file exclusion patterns work correctly
// when there are similarly named files and directories.
func TestGetAllFilePathsExcludeFileVsDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a file and a directory with similar names
	filePath := filepath.Join(rootDir, "build")
	err := os.WriteFile(filePath, []byte("file content"), 0644)
	require.NoError(t, err, "Failed to create file")

	dirPath := filepath.Join(rootDir, "build_dir")
	err = os.Mkdir(dirPath, 0755)
	require.NoError(t, err, "Failed to create directory")

	fileInDir := filepath.Join(dirPath, "file.txt")
	err = os.WriteFile(fileInDir, []byte("dir file content"), 0644)
	require.NoError(t, err, "Failed to create file in directory")

	// Exclude "build" which is a file
	excludePatterns := []string{"build"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	// Should exclude the file but include the directory and its contents
	expected := []string{
		dirPath,
		fileInDir,
	}
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsExcludeHiddenDirectory tests that hidden directories (like .git)
// can be properly excluded along with their contents.
func TestGetAllFilePathsExcludeHiddenDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a hidden directory and a file inside it
	hiddenDir := filepath.Join(rootDir, ".git")
	err := os.Mkdir(hiddenDir, 0755)
	require.NoError(t, err, "Failed to create hidden directory")

	configFile := filepath.Join(hiddenDir, "config")
	err = os.WriteFile(configFile, []byte("config content"), 0644)
	require.NoError(t, err, "Failed to create config file")

	// Exclude ".git/" directory
	excludePatterns := []string{".git/"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.Empty(t, filePaths, "Expected no files when excluding hidden directory")
}

// TestGetAllFilePathsIncludeExcludeOverlap tests the interaction between include and exclude patterns,
// ensuring that exclude patterns take precedence over include patterns.
func TestGetAllFilePathsIncludeExcludeOverlap(t *testing.T) {
	rootDir := t.TempDir()

	// Create two .go files
	file1 := filepath.Join(rootDir, "file1.go")
	err := os.WriteFile(file1, []byte("content1"), 0644)
	require.NoError(t, err, "Failed to create file1")

	file2 := filepath.Join(rootDir, "file2.go")
	err = os.WriteFile(file2, []byte("content2"), 0644)
	require.NoError(t, err, "Failed to create file2")

	// Include all .go files, but exclude file2.go
	includePatterns := []string{"**/*.go"}
	excludePatterns := []string{"file2.go"}
	expected := []string{file1}

	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsCaseSensitivity tests that file pattern matching is case-sensitive.
func TestGetAllFilePathsCaseSensitivity(t *testing.T) {
	rootDir := t.TempDir()

	// Create files with different cases
	file1 := filepath.Join(rootDir, "README_upper")
	err := os.WriteFile(file1, []byte("uppercase"), 0644)
	require.NoError(t, err, "Failed to create uppercase file")

	file2 := filepath.Join(rootDir, "readme_lower")
	err = os.WriteFile(file2, []byte("lowercase"), 0644)
	require.NoError(t, err, "Failed to create lowercase file")

	// Exclude "README_upper"
	excludePatterns := []string{"README_upper"}
	expected := []string{file2}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsExcludeNonExistingDirectory tests that excluding non-existent directories
// does not affect the inclusion of existing files.
func TestGetAllFilePathsExcludeNonExistingDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a test file
	filePath := filepath.Join(rootDir, "file.txt")
	err := os.WriteFile(filePath, []byte("content"), 0644)
	require.NoError(t, err, "Failed to create file")

	// Exclude a non-existing directory
	excludePatterns := []string{"nonexistent_dir/"}
	expected := []string{filePath}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsExcludeEmptyPattern tests that empty exclude patterns are properly handled
// and do not affect file inclusion.
func TestGetAllFilePathsExcludeEmptyPattern(t *testing.T) {
	rootDir := t.TempDir()

	// Create a test file
	filePath := filepath.Join(rootDir, "file.txt")
	err := os.WriteFile(filePath, []byte("content"), 0644)
	require.NoError(t, err, "Failed to create file")

	// Test with an empty exclude pattern
	excludePatterns := []string{""}
	expected := []string{filePath}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsExcludeSymlink tests that symbolic links can be properly excluded
// while preserving access to the original files.
func TestGetAllFilePathsExcludeSymlink(t *testing.T) {
	rootDir := t.TempDir()

	// Create a target directory and file
	targetDir := filepath.Join(rootDir, "target")
	err := os.Mkdir(targetDir, 0755)
	require.NoError(t, err, "Failed to create target directory")

	targetFile := filepath.Join(targetDir, "file.txt")
	err = os.WriteFile(targetFile, []byte("content"), 0644)
	require.NoError(t, err, "Failed to create target file")

	// Create a symlink to the directory
	symlinkDir := filepath.Join(rootDir, "symlink")
	err = os.Symlink(targetDir, symlinkDir)
	require.NoError(t, err, "Failed to create symlink")

	// Exclude the symlink
	excludePatterns := []string{"symlink/"}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	// Should only include the target directory and its file
	expected := []string{
		targetDir,
		targetFile,
	}
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}
