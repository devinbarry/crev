package files_test

import (
	"github.com/devinbarry/crev/internal/files"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

// TestGetAllFilePaths tests the basic functionality to get all file paths starting from a root path,
// including subdirectories and their contents.
func TestGetAllFilePaths(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"subdir/file1.txt": "content1",
		"subdir/file2.txt": "content2",
	}
	createFiles(t, rootDir, fileStructure)

	includePatterns := []string{"**/*"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, nil, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{
		filepath.Join(rootDir, "subdir"),
		filepath.Join(rootDir, "subdir/file1.txt"),
		filepath.Join(rootDir, "subdir/file2.txt"),
	}
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsEmpty tests all args empty
func TestGetAllFilePathsEmpty(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"subdir/file1.txt": "content1",
		"subdir/file2.txt": "content2",
	}
	createFiles(t, rootDir, fileStructure)

	filePaths, err := files.GetAllFilePaths(rootDir, nil, nil, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	// We should get nothing at all included
	require.ElementsMatch(t, nil, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsWithExcludePattern tests the functionality of exclude patterns with globbing,
// ensuring that excluded directories and their contents are properly filtered out.
func TestGetAllFilePathsWithExcludePattern(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file1.go":          "content1",
		"subdir_1/file2.go": "content2",
		"subdir_1/file4.go": "content2",
		"subdir_2/file3.go": "content3",
	}
	createFiles(t, rootDir, fileStructure)

	// Exclude subdir_1 and its contents using glob pattern
	includePatterns := []string{"**/*"}
	excludePatterns := []string{"subdir_1/**"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "subdir_2"),
		filepath.Join(rootDir, "subdir_2", "file3.go"),
	}
	notExpected := []string{
		filepath.Join(rootDir, "subdir_1/file2.go"),
		filepath.Join(rootDir, "subdir_1/file4.go"),
	}
	assertFileSetMatches(t, filePaths, expected, notExpected,
		"Incorrect paths returned")
}

// TestGetAllFilePathsWithIncludePattern tests the functionality of include patterns with globbing,
// ensuring that only files matching the include pattern are returned.
func TestGetAllFilePathsWithIncludePattern(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file1.go":           "content1",
		"file2.txt":          "content2",
		"subdir_1/file3.go":  "content3",
		"subdir_1/file4.txt": "content4",
		"subdir_2/file5.go":  "content5",
	}
	createFiles(t, rootDir, fileStructure)

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "subdir_1", "file3.go"),
		filepath.Join(rootDir, "subdir_2", "file5.go"),
	}

	// Include only .go files using glob pattern
	includePatterns := []string{"**/*.go"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, nil, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsIncludeAndExcludePatterns tests combining include and exclude patterns,
// verifying that both patterns work together correctly.
func TestGetAllFilePathsIncludeAndExcludePatterns(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file1.go":           "content1",
		"file2.go":           "content2",
		"subdir_1/file3.go":  "content3",
		"subdir_1/file4.txt": "content4",
		"subdir_2/file5.go":  "content5",
		"subdir_2/file6.go":  "content6",
	}
	createFiles(t, rootDir, fileStructure)

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "file2.go"),
		filepath.Join(rootDir, "subdir_1", "file3.go"),
	}

	// Include all .go files but exclude subdir_2
	includePatterns := []string{"**/*.go"}
	excludePatterns := []string{"subdir_2/**"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsWithExtensionExcludePatterns tests excluding files by extension using glob patterns,
// verifying that multiple extension exclusions work correctly across nested directories.
func TestGetAllFilePathsWithExtensionExcludePatterns(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file1.go":                           "content1",
		"subdir_1/file2.go":                  "content2",
		"subdir_1/nested_subdir_1/file3.go":  "content3",
		"subdir_2/file4.txt":                 "content4",
		"subdir_2/nested_subdir_2/file5.md":  "content5",
		"subdir_2/nested_subdir_2/file6.txt": "content6",
	}
	createFiles(t, rootDir, fileStructure)

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "subdir_1", "file2.go"),
		filepath.Join(rootDir, "subdir_1", "nested_subdir_1", "file3.go"),
	}

	// Exclude .txt and .md files using glob patterns
	includePatterns := []string{"**/*.go"}
	excludePatterns := []string{"**/*.txt", "**/*.md"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsWithExtensionExcludePatterns tests excluding files by extension using glob patterns,
// verifying that multiple extension exclusions work correctly across nested directories.
func TestGetAllFilePathsIncludeSubdirectories(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file1.go":                           "content1",
		"subdir_1/file2.go":                  "content2",
		"subdir_1/nested_subdir_1/file3.go":  "content3",
		"subdir_2/file4.txt":                 "content4",
		"subdir_2/nested_subdir_2/file5.md":  "content5",
		"subdir_2/nested_subdir_2/file6.txt": "content6",
	}
	createFiles(t, rootDir, fileStructure)

	// This test is like the one above, but here we have a much wider glob include pattern
	// This patterns means we get subdirectories included in the output
	includePatterns := []string{"**/*"}
	excludePatterns := []string{"**/*.txt", "**/*.md"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "subdir_1"),
		filepath.Join(rootDir, "subdir_1", "file2.go"),
		filepath.Join(rootDir, "subdir_1", "nested_subdir_1"),
		filepath.Join(rootDir, "subdir_1", "nested_subdir_1", "file3.go"),
	}
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsWithExplicitFiles tests including explicit files regardless of patterns,
// verifying that explicitly included files override exclude patterns.
func TestGetAllFilePathsWithExplicitFiles(t *testing.T) {
	rootDir := t.TempDir()

	createFile(t, filepath.Join(rootDir, "file1.go"), "content1")
	createFile(t, filepath.Join(rootDir, "file2.txt"), "content2")

	// Exclude all .txt files but include file2.txt explicitly
	excludePatterns := []string{"**/*.txt"}
	explicitFiles := []string{filepath.Join(rootDir, "file2.txt")}
	expected := []string{
		filepath.Join(rootDir, "file2.txt"),
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, explicitFiles)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}
