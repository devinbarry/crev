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

	fileStructure := map[string]string{
		"dir/file.txt": "content",
	}
	createFiles(t, rootDir, fileStructure)

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

	fileStructure := map[string]string{
		"build":              "file content",
		"build_dir/file.txt": "dir file content",
	}
	createFiles(t, rootDir, fileStructure)

	// Exclude "build" which is a file
	excludePatterns := []string{"build"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{
		filepath.Join(rootDir, "build_dir"),
		filepath.Join(rootDir, "build_dir/file.txt"),
	}
	notExpected := []string{
		filepath.Join(rootDir, "build"),
	}

	assertFileSetMatches(t, filePaths, expected, notExpected,
		"Should exclude 'build' file but include 'build_dir' directory and its contents")
}

// TestGetAllFilePathsExcludeHiddenDirectory tests that hidden directories (like .git)
// can be properly excluded along with their contents.
func TestGetAllFilePathsExcludeHiddenDirectory(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"tests/format_test.go":   "tests",
		"tests/globbing_test.go": "tests",
		".git/config":            "config content",
		".git/FETCH_HEAD":        "HEAD content",
		".git/COMMIT":            "COMMIT content",
	}
	createFiles(t, rootDir, fileStructure)

	includePatterns := []string{"**/*"}
	excludePatterns := []string{".git/"} // Exclude ".git/" directory
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{
		filepath.Join(rootDir, "tests/"),
		filepath.Join(rootDir, "tests/format_test.go"),
		filepath.Join(rootDir, "tests/globbing_test.go"),
	}
	notExpected := []string{
		filepath.Join(rootDir, ".git/"),
		filepath.Join(rootDir, ".git/config"),
		filepath.Join(rootDir, ".git/FETCH_HEAD"),
		filepath.Join(rootDir, ".git/COMMIT"),
	}
	assertFileSetMatches(t, filePaths, expected, notExpected,
		"Should exclude hidden directory but include repo files")
}

// TestGetAllFilePathsIncludeExcludeOverlap tests the interaction between include and exclude patterns,
// ensuring that exclude patterns take precedence over include patterns.
func TestGetAllFilePathsIncludeExcludeOverlap(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file1.go": "content1",
		"file2.go": "content2",
	}
	createFiles(t, rootDir, fileStructure)

	// Include all .go files, but exclude file2.go
	includePatterns := []string{"**/*.go"}
	excludePatterns := []string{"file2.go"}
	expected := []string{filepath.Join(rootDir, "file1.go")}

	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsCaseSensitivity tests that file pattern matching is case-sensitive.
func TestGetAllFilePathsCaseSensitivity(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"README_upper": "uppercase",
		"readme_lower": "lowercase",
	}
	createFiles(t, rootDir, fileStructure)

	// Exclude "README_upper"
	excludePatterns := []string{"README_upper"}
	expected := []string{filepath.Join(rootDir, "readme_lower")}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsExcludeNonExistingDirectory tests that excluding non-existent directories
// does not affect the inclusion of existing files.
func TestGetAllFilePathsExcludeNonExistingDirectory(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file.txt": "content",
	}
	createFiles(t, rootDir, fileStructure)

	// Exclude a non-existing directory
	excludePatterns := []string{"nonexistent_dir/"}
	expected := []string{filepath.Join(rootDir, "file.txt")}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsExcludeEmptyPattern tests that empty exclude patterns are properly handled
// and do not affect file inclusion.
func TestGetAllFilePathsExcludeEmptyPattern(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file.txt": "content",
	}
	createFiles(t, rootDir, fileStructure)

	// Test with an empty exclude pattern
	excludePatterns := []string{""}
	expected := []string{filepath.Join(rootDir, "file.txt")}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsExcludeSymlink tests that symbolic links can be properly excluded
// while preserving access to the original files.
func TestGetAllFilePathsExcludeSymlink(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"target/file.txt": "content",
	}
	createFiles(t, rootDir, fileStructure)

	// Create a symlink to the directory
	targetDir := filepath.Join(rootDir, "target")
	symlinkDir := filepath.Join(rootDir, "symlink")
	err := os.Symlink(targetDir, symlinkDir)
	require.NoError(t, err, "Failed to create symlink")

	// Exclude the symlink
	excludePatterns := []string{"symlink/"}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	// Should only include the target directory and its file
	expected := []string{
		targetDir,
		filepath.Join(rootDir, "target/file.txt"),
	}
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}
