package files

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

// TestGetAllFilePathsExcludeDirTrailingSlash tests that directories are correctly excluded
// whether the exclude pattern has a trailing slash or not.
func TestGetAllFilePathsExcludeDirTrailingSlash(t *testing.T) {
	// Setup temporary directory
	rootDir := t.TempDir()

	// Create a directory and a file inside it
	dirPath := filepath.Join(rootDir, "dir")
	createDir(t, dirPath)
	fileInDir := filepath.Join(dirPath, "file.txt")
	createFile(t, fileInDir, "content")

	// Test excluding directory without trailing slash
	excludePatterns := []string{"dir"}
	filePaths, err := GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")
	require.Len(t, filePaths, 0, "Expected 0 files after exclusion")

	// Test excluding directory with trailing slash
	excludePatterns = []string{"dir/"}
	filePaths, err = GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")
	require.Len(t, filePaths, 0, "Expected 0 files after exclusion with trailing slash")
}

// TestGetAllFilePathsExcludeFileVsDirectory tests excluding a file vs a directory with the same name.
func TestGetAllFilePathsExcludeFileVsDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a file named "build"
	filePath := filepath.Join(rootDir, "build")
	createFile(t, filePath, "file content")

	// Create a directory named "build_dir" with a file inside
	dirPath := filepath.Join(rootDir, "build_dir")
	createDir(t, dirPath)
	fileInDir := filepath.Join(dirPath, "file.txt")
	createFile(t, fileInDir, "dir file content")

	// Exclude "build" which is a file
	excludePatterns := []string{"build"}
	filePaths, err := GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")

	// Expected: Only the directory and its contents should be included
	expected := []string{
		dirPath,
		fileInDir,
	}
	require.ElementsMatch(t, expected, filePaths, "Excluded file should be omitted, directory should remain")
}

// TestGetAllFilePathsExcludeHiddenDirectory tests excluding hidden directories like ".git".
func TestGetAllFilePathsExcludeHiddenDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a hidden directory and a file inside it
	hiddenDir := filepath.Join(rootDir, ".git")
	createDir(t, hiddenDir)
	configFile := filepath.Join(hiddenDir, "config")
	createFile(t, configFile, "[core]")

	// Exclude the ".git" directory
	excludePatterns := []string{".git/"}
	filePaths, err := GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")
	require.Len(t, filePaths, 0, "Expected 0 files after excluding hidden directory")
}

// TestGetAllFilePathsIncludeExcludeOverlap tests that include and exclude patterns interact correctly.
func TestGetAllFilePathsIncludeExcludeOverlap(t *testing.T) {
	rootDir := t.TempDir()

	// Create two .go files
	file1 := filepath.Join(rootDir, "file1.go")
	createFile(t, file1, "package main")
	file2 := filepath.Join(rootDir, "file2.go")
	createFile(t, file2, "package main")

	// Exclude "file2.go" while including all .go files
	includePatterns := []string{"**/*.go"}
	excludePatterns := []string{"file2.go"}
	filePaths, err := GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")

	// Only "file1.go" should be included
	expected := []string{file1}
	require.ElementsMatch(t, expected, filePaths, "Only file1.go should be included after exclusion")
}

// TestGetAllFilePathsCaseSensitivity tests that file exclusion is case-sensitive.
func TestGetAllFilePathsCaseSensitivity(t *testing.T) {
	rootDir := t.TempDir()

	// Create files with different cases
	file1 := filepath.Join(rootDir, "README_upper")
	createFile(t, file1, "uppercase")
	file2 := filepath.Join(rootDir, "readme_lower")
	createFile(t, file2, "lowercase")

	// Exclude "README_upper"
	excludePatterns := []string{"README_upper"}
	filePaths, err := GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")

	// Only "readme_lower" should be included
	expected := []string{file2}
	require.ElementsMatch(t, expected, filePaths, "Only readme_lower should be included after excluding README_upper")
}

// TestGetAllFilePathsExcludeNonExistingDirectory tests that excluding a non-existing directory doesn't affect existing files.
func TestGetAllFilePathsExcludeNonExistingDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a file
	filePath := filepath.Join(rootDir, "file.txt")
	createFile(t, filePath, "content")

	// Exclude a non-existing directory
	excludePatterns := []string{"nonexistent_dir/"}
	filePaths, err := GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")

	// The existing file should still be included
	expected := []string{filePath}
	require.ElementsMatch(t, expected, filePaths, "Existing files should remain unaffected by excluding non-existent directories")
}

// TestGetAllFilePathsExcludeEmptyPattern tests that empty exclude patterns are ignored.
func TestGetAllFilePathsExcludeEmptyPattern(t *testing.T) {
	rootDir := t.TempDir()

	// Create a file
	filePath := filepath.Join(rootDir, "file.txt")
	createFile(t, filePath, "content")

	// Exclude patterns list contains an empty string
	excludePatterns := []string{""}
	filePaths, err := GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")

	// The existing file should still be included
	expected := []string{filePath}
	require.ElementsMatch(t, expected, filePaths, "Empty exclude patterns should be ignored")
}

// TestGetAllFilePathsExcludeSymlink tests that symlinks can be excluded properly.
func TestGetAllFilePathsExcludeSymlink(t *testing.T) {
	rootDir := t.TempDir()

	// Create a target directory and a file inside it
	targetDir := filepath.Join(rootDir, "target")
	createDir(t, targetDir)
	targetFile := filepath.Join(targetDir, "file.txt")
	createFile(t, targetFile, "content")

	// Create a symlink to the target directory
	symlinkDir := filepath.Join(rootDir, "symlink")
	err := os.Symlink(targetDir, symlinkDir)
	require.NoError(t, err, "Failed to create symlink")

	// Exclude the symlink directory
	excludePatterns := []string{"symlink/"}
	filePaths, err := GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "Failed to get all file paths")

	// The symlink and its contents should be excluded
	expected := []string{targetDir, targetFile}
	require.ElementsMatch(t, expected, filePaths, "Symlinked directories should be excluded correctly")
}
