package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devinbarry/crev/internal/files"
)

// TestGetAllFilePaths tests the basic functionality to get all file paths starting from a root path.
func TestGetAllFilePaths(t *testing.T) {
	rootDir := t.TempDir()

	subDir := filepath.Join(rootDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(rootDir, "file1.txt"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{
		filepath.Join(rootDir, "file1.txt"),
		subDir,
		filepath.Join(subDir, "file2.txt"),
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	for _, exp := range expected {
		found := false
		for _, fp := range filePaths {
			if fp == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %s not found in result", exp)
		}
	}
}

// TestGetAllFilePathsWithExcludePattern tests the functionality of exclude patterns with globbing.
func TestGetAllFilePathsWithExcludePattern(t *testing.T) {
	rootDir := t.TempDir()

	subDir1 := filepath.Join(rootDir, "subdir_1")
	err := os.Mkdir(subDir1, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	subDir2 := filepath.Join(rootDir, "subdir_2")
	err = os.Mkdir(subDir2, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(rootDir, "file1.go"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir1, "file2.go"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir2, "file3.go"), []byte("content3"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		subDir2,
		filepath.Join(subDir2, "file3.go"),
	}

	// Exclude subdir_1 and its contents using glob pattern
	excludePatterns := []string{"subdir_1/**"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	for _, exp := range expected {
		found := false
		for _, fp := range filePaths {
			if fp == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %s not found in result", exp)
		}
	}
}

// TestGetAllFilePathsWithIncludePattern tests the functionality of include patterns with globbing.
func TestGetAllFilePathsWithIncludePattern(t *testing.T) {
	rootDir := t.TempDir()

	subDir1 := filepath.Join(rootDir, "subdir_1")
	err := os.Mkdir(subDir1, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	subDir2 := filepath.Join(rootDir, "subdir_2")
	err = os.Mkdir(subDir2, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(rootDir, "file1.go"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(rootDir, "file2.txt"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir1, "file3.go"), []byte("content3"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir1, "file4.txt"), []byte("content4"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir2, "file5.go"), []byte("content5"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(subDir1, "file3.go"),
		filepath.Join(subDir2, "file5.go"),
	}

	// Include only .go files using glob pattern
	includePatterns := []string{"**/*.go"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	for _, exp := range expected {
		found := false
		for _, fp := range filePaths {
			if fp == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %s not found in result", exp)
		}
	}
}

// TestGetAllFilePathsIncludeAndExcludePatterns tests combining include and exclude patterns.
func TestGetAllFilePathsIncludeAndExcludePatterns(t *testing.T) {
	rootDir := t.TempDir()

	subDir1 := filepath.Join(rootDir, "subdir_1")
	err := os.Mkdir(subDir1, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	subDir2 := filepath.Join(rootDir, "subdir_2")
	err = os.Mkdir(subDir2, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(rootDir, "file1.go"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(rootDir, "file2.go"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir1, "file3.go"), []byte("content3"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir1, "file4.txt"), []byte("content4"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir2, "file5.go"), []byte("content5"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir2, "file6.go"), []byte("content6"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "file2.go"),
		filepath.Join(subDir1, "file3.go"),
	}

	// Include all .go files but exclude subdir_2
	includePatterns := []string{"**/*.go"}
	excludePatterns := []string{"subdir_2/**"}
	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	for _, exp := range expected {
		found := false
		for _, fp := range filePaths {
			if fp == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %s not found in result", exp)
		}
	}
}

// TestGetAllFilePathsWithExtensionExcludePatterns tests excluding files by extension using glob patterns.
func TestGetAllFilePathsWithExtensionExcludePatterns(t *testing.T) {
	rootDir := t.TempDir()

	// Create subdirectories and nested subdirectories
	subDir1 := filepath.Join(rootDir, "subdir_1")
	err := os.Mkdir(subDir1, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	nestedSubDir1 := filepath.Join(subDir1, "nested_subdir_1")
	err = os.Mkdir(nestedSubDir1, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	subDir2 := filepath.Join(rootDir, "subdir_2")
	err = os.Mkdir(subDir2, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	nestedSubDir2 := filepath.Join(subDir2, "nested_subdir_2")
	err = os.Mkdir(nestedSubDir2, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Create files in various directories
	err = os.WriteFile(filepath.Join(rootDir, "file1.go"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir1, "file2.go"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(nestedSubDir1, "file3.go"), []byte("content3"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir2, "file4.txt"), []byte("content4"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(nestedSubDir2, "file5.md"), []byte("content5"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(nestedSubDir2, "file6.txt"), []byte("content6"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Expected result: exclude .txt and .md files, keep the rest
	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(subDir1, "file2.go"),
		filepath.Join(nestedSubDir1, "file3.go"),
		subDir1,
		nestedSubDir1,
		subDir2,
		nestedSubDir2,
	}

	// Exclude .txt and .md files using glob patterns
	excludePatterns := []string{"**/*.txt", "**/*.md"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the number of files found
	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	// Verify that each expected file is in the result
	for _, exp := range expected {
		found := false
		for _, fp := range filePaths {
			if fp == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %s not found in result", exp)
		}
	}
}

// TestGetAllFilePathsWithExplicitFiles tests including explicit files regardless of patterns.
func TestGetAllFilePathsWithExplicitFiles(t *testing.T) {
	rootDir := t.TempDir()

	// Create files
	err := os.WriteFile(filepath.Join(rootDir, "file1.go"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(rootDir, "file2.txt"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exclude all .txt files but include file2.txt explicitly
	excludePatterns := []string{"**/*.txt"}
	explicitFiles := []string{filepath.Join(rootDir, "file2.txt")}
	expected := []string{
		filepath.Join(rootDir, "file1.go"),
		filepath.Join(rootDir, "file2.txt"),
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, explicitFiles)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	for _, exp := range expected {
		found := false
		for _, fp := range filePaths {
			if fp == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %s not found in result", exp)
		}
	}
}

// TestGetContentMapOfFiles tests reading the content of files and handling empty directories.
func TestGetContentMapOfFiles(t *testing.T) {
	rootDir := t.TempDir()

	subDir1 := filepath.Join(rootDir, "subdir_1")
	subDir2 := filepath.Join(rootDir, "subdir_2")
	err := os.Mkdir(subDir1, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.Mkdir(subDir2, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(rootDir, "file1.txt"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir1, "file2.txt"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	filePaths := []string{
		filepath.Join(rootDir, "file1.txt"),
		subDir1,
		filepath.Join(subDir1, "file2.txt"),
		subDir2,
	}

	fileContentMap, err := files.GetContentMapOfFiles(filePaths, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fileContentMap[filepath.Join(rootDir, "file1.txt")] != "content1" {
		t.Errorf("expected content1, got %s", fileContentMap[filepath.Join(rootDir, "file1.txt")])
	}

	if fileContentMap[filepath.Join(subDir1, "file2.txt")] != "content2" {
		t.Errorf("expected content2, got %s", fileContentMap[filepath.Join(subDir1, "file2.txt")])
	}

	// subDir1 is not empty, so it should not be present in the map
	if _, ok := fileContentMap[subDir1]; ok {
		t.Errorf("directory with files should not be present %s", fileContentMap[subDir1])
	}

	// subDir2 is empty, so it should be present in the map with "empty directory"
	if fileContentMap[subDir2] != "empty directory" {
		t.Errorf("expected empty directory, got %s", fileContentMap[subDir2])
	}
}
