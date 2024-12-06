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

	subDir := filepath.Join(rootDir, "subdir")
	createDir(t, subDir)
	createFile(t, filepath.Join(rootDir, "file1.txt"), "content1")
	createFile(t, filepath.Join(subDir, "file2.txt"), "content2")

	expected := []string{
		filepath.Join(rootDir, "file1.txt"),
		subDir,
		filepath.Join(subDir, "file2.txt"),
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, nil, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

// TestGetAllFilePathsWithExcludePattern tests the functionality of exclude patterns with globbing,
// ensuring that excluded directories and their contents are properly filtered out.
func TestGetAllFilePathsWithExcludePattern(t *testing.T) {
	rootDir := t.TempDir()

	fileStructure := map[string]string{
		"file1.go":          "content1",
		"subdir_1/file2.go": "content2",
		"subdir_2/file3.go": "content3",
	}
	createFiles(t, rootDir, fileStructure)

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "subdir_2"),
		filepath.Join(rootDir, "subdir_2", "file3.go"),
	}

	// Exclude subdir_1 and its contents using glob pattern
	excludePatterns := []string{"subdir_1/**"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
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
		filepath.Join(rootDir, "subdir_1"),
		filepath.Join(rootDir, "subdir_1", "file2.go"),
		filepath.Join(rootDir, "subdir_1", "nested_subdir_1"),
		filepath.Join(rootDir, "subdir_1", "nested_subdir_1", "file3.go"),
		filepath.Join(rootDir, "subdir_2"),
		filepath.Join(rootDir, "subdir_2", "nested_subdir_2"),
	}

	// Exclude .txt and .md files using glob patterns
	excludePatterns := []string{"**/*.txt", "**/*.md"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	require.NoError(t, err, "GetAllFilePaths failed")
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
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "file2.txt"),
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, explicitFiles)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expected, filePaths, "Incorrect paths returned")
}

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
