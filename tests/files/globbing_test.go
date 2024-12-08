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

// TestGetAllFilePathsExcludeFileVsDirectory tests excluding a file vs a directory with the same name.
func TestGetAllFilePathsExcludeFileVsDirectory(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"build":              "file content",
		"build_dir/file.txt": "dir file content",
	}
	createFiles(t, rootDir, fileStructure)

	// Exclude "build" which is a file
	includePatterns := []string{"**/*"}
	excludePatterns := []string{"build"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{
		"build_dir",
		"build_dir/file.txt",
	}
	notExpected := []string{
		"build",
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
		"tests",
		"tests/format_test.go",
		"tests/globbing_test.go",
	}
	notExpected := []string{
		".git",
		".git/config",
		".git/FETCH_HEAD",
		".git/COMMIT",
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
		"file3.go": "content3",
	}
	createFiles(t, rootDir, fileStructure)

	// Include all .go files, but exclude file2.go
	includePatterns := []string{"**/*.go"}
	excludePatterns := []string{"file2.go"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	// file1 and file3 but not file2
	expected := []string{
		"file1.go",
		"file3.go",
	}
	notExpected := []string{"file2.go"}
	assertFileSetMatches(t, filePaths, expected, notExpected,
		"Only file1.go should be included after exclusion")
}

// TestGetAllFilePathsCaseSensitivity tests that file pattern matching is case-sensitive.
func TestGetAllFilePathsCaseSensitivity(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"README": "uppercase",
		"readme": "lowercase",
	}
	createFiles(t, rootDir, fileStructure)

	// Exclude "README uppercase only"
	includePatterns := []string{"**/*"}
	excludePatterns := []string{"README"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{"readme"}
	notExpected := []string{"README"}
	assertFileSetMatches(t, filePaths, expected, notExpected,
		"Only readme should be included after excluding README")
}

// TestGetAllFilePathsExcludeNonExistingDirectory tests that excluding a non-existing directory doesn't affect existing files.
func TestGetAllFilePathsExcludeNonExistingDirectory(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file.txt": "content",
		"money.py": "# time to get rich",
	}
	createFiles(t, rootDir, fileStructure)

	// Exclude a non-existing directory
	includePatterns := []string{"**/*"}
	excludePatterns := []string{"nonexistent_dir/"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{
		"file.txt",
		"money.py",
	}
	assertFileSetMatches(t, filePaths, expected, nil,
		"Existing files should remain unaffected by excluding non-existent directories")
}

// TestGetAllFilePathsExcludeEmptyPattern tests that empty exclude patterns are ignored.
func TestGetAllFilePathsExcludeEmptyPattern(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file.txt": "content",
	}
	createFiles(t, rootDir, fileStructure)

	// Exclude patterns list contains an empty string
	includePatterns := []string{"**/*"}
	excludePatterns := []string{""}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{"file.txt"}
	assertFileSetMatches(t, filePaths, expected, nil,
		"Empty exclude patterns should be ignored")
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
	includePatterns := []string{"**/*"}
	excludePatterns := []string{"symlink/"}

	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	// Should only include the target directory and its file
	expected := []string{
		"target",
		"target/file.txt",
	}
	// The symlink and its contents should be excluded
	notExpected := []string{
		"symlink",
		"symlink/file.txt",
	}
	assertFileSetMatches(t, filePaths, expected, notExpected,
		"Symlinked directories should be excluded correctly")
}
