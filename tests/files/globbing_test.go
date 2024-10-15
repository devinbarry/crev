package files_test

import (
	"github.com/devinbarry/crev/internal/files"
	"os"
	"path/filepath"
	"testing"
)

func TestGetAllFilePathsExcludeDirTrailingSlash(t *testing.T) {
	rootDir := t.TempDir()

	// Create directories and files
	dirPath := filepath.Join(rootDir, "dir")
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	fileInDir := filepath.Join(dirPath, "file.txt")
	err = os.WriteFile(fileInDir, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exclude pattern without trailing slash
	excludePatterns := []string{"dir"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != 0 {
		t.Errorf("expected 0 files, got %d", len(filePaths))
	}

	// Exclude pattern with trailing slash
	excludePatterns = []string{"dir/"}
	filePaths, err = files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != 0 {
		t.Errorf("expected 0 files, got %d", len(filePaths))
	}
}

func TestGetAllFilePathsExcludeFileVsDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a file and a directory with the same name
	filePath := filepath.Join(rootDir, "build")
	err := os.WriteFile(filePath, []byte("file content"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dirPath := filepath.Join(rootDir, "build_dir")
	err = os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	fileInDir := filepath.Join(dirPath, "file.txt")
	err = os.WriteFile(fileInDir, []byte("dir file content"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exclude "build" which is a file
	excludePatterns := []string{"build"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should exclude the file but include the directory and its contents
	expected := []string{
		dirPath,
		fileInDir,
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

func TestGetAllFilePathsExcludeHiddenDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a hidden directory and a file inside it
	hiddenDir := filepath.Join(rootDir, ".git")
	err := os.Mkdir(hiddenDir, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	configFile := filepath.Join(hiddenDir, "config")
	err = os.WriteFile(configFile, []byte("config content"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exclude ".git/" directory
	excludePatterns := []string{".git/"}
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should not include the .git directory or its contents
	if len(filePaths) != 0 {
		t.Errorf("expected 0 files, got %d", len(filePaths))
	}
}

func TestGetAllFilePathsIncludeExcludeOverlap(t *testing.T) {
	rootDir := t.TempDir()

	// Create files
	file1 := filepath.Join(rootDir, "file1.go")
	err := os.WriteFile(file1, []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	file2 := filepath.Join(rootDir, "file2.go")
	err = os.WriteFile(file2, []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Include all .go files, but exclude file2.go
	includePatterns := []string{"**/*.go"}
	excludePatterns := []string{"file2.go"}
	expected := []string{
		file1,
	}

	filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	if filePaths[0] != expected[0] {
		t.Errorf("expected path %s, got %s", expected[0], filePaths[0])
	}
}

func TestGetAllFilePathsCaseSensitivity(t *testing.T) {
	rootDir := t.TempDir()

	// Create files with different names
	file1 := filepath.Join(rootDir, "README_upper")
	err := os.WriteFile(file1, []byte("uppercase"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	file2 := filepath.Join(rootDir, "readme_lower")
	err = os.WriteFile(file2, []byte("lowercase"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exclude "README_upper"
	excludePatterns := []string{"README_upper"}
	expected := []string{
		file2,
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	if filePaths[0] != expected[0] {
		t.Errorf("expected path %s, got %s", expected[0], filePaths[0])
	}
}

func TestGetAllFilePathsExcludeNonExistingDirectory(t *testing.T) {
	rootDir := t.TempDir()

	// Create a file
	filePath := filepath.Join(rootDir, "file.txt")
	err := os.WriteFile(filePath, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exclude a non-existing directory
	excludePatterns := []string{"nonexistent_dir/"}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// The existing file should still be included
	expected := []string{
		filePath,
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	if filePaths[0] != expected[0] {
		t.Errorf("expected path %s, got %s", expected[0], filePaths[0])
	}
}

func TestGetAllFilePathsExcludeEmptyPattern(t *testing.T) {
	rootDir := t.TempDir()

	// Create a file
	filePath := filepath.Join(rootDir, "file.txt")
	err := os.WriteFile(filePath, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exclude patterns list contains an empty string
	excludePatterns := []string{""}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// The existing file should still be included
	expected := []string{
		filePath,
	}

	if len(filePaths) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(filePaths))
	}

	if filePaths[0] != expected[0] {
		t.Errorf("expected path %s, got %s", expected[0], filePaths[0])
	}
}

func TestGetAllFilePathsExcludeSymlink(t *testing.T) {
	rootDir := t.TempDir()

	// Create a directory and a file
	targetDir := filepath.Join(rootDir, "target")
	err := os.Mkdir(targetDir, 0755)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	targetFile := filepath.Join(targetDir, "file.txt")
	err = os.WriteFile(targetFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Create a symlink to the directory
	symlinkDir := filepath.Join(rootDir, "symlink")
	err = os.Symlink(targetDir, symlinkDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exclude the symlink
	excludePatterns := []string{"symlink/"}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should only include the target directory and its file
	expected := []string{
		targetDir,
		targetFile,
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
