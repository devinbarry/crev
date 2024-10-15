package formatting_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devinbarry/crev/internal/formatting"
)

// TestCreateProjectString tests the combined output of directory tree and file contents.
func TestCreateProjectString(t *testing.T) {
	projectTree, err := os.ReadFile("../test_data/expected_path_tree_full.txt")
	if err != nil {
		t.Fatalf("error reading test file: %v", err)
	}
	fileContentMap := map[string]string{
		"cmd/ai-code-review/main.go":    "package main\n",
		"internal/files/filtering.go":   "package files\n",
		"internal/formatting/format.go": "package formatting\n",
		"go.mod":                        "go mod\n",
	}
	expectedProjectString, err := os.ReadFile("../test_data/expected_project_output.txt")
	if err != nil {
		t.Fatalf("error reading test file: %v", err)
	}
	expected := string(expectedProjectString)
	result := formatting.CreateProjectString(string(projectTree), fileContentMap)
	expectedStr := strings.TrimSpace(string(expected))
	resultStr := strings.TrimSpace(result)
	if resultStr != expectedStr {
		t.Errorf("CreateProjectString: expected \n%s\n, got \n%s\n", expectedStr, resultStr)
	}
}

// TestGeneratePathTree tests the directory tree generation independently.
func TestGeneratePathTree(t *testing.T) {
	paths := []string{
		"cmd",
		"cmd/ai-code-review",
		"cmd/ai-code-review/main.go",
		"internal",
		"internal/files",
		"internal/files/filtering.go",
		"internal/formatting",
		"internal/formatting/format.go",
		"go.mod",
	}
	expectedTree, err := os.ReadFile("../test_data/expected_path_tree_full.txt")
	if err != nil {
		t.Fatalf("error reading test file: %v", err)
	}

	resultTree := formatting.GeneratePathTree(paths)
	expectedStr := strings.TrimSpace(string(expectedTree))
	resultStr := strings.TrimSpace(resultTree)

	if resultStr != expectedStr {
		t.Errorf("GeneratePathTree: expected \n%s\n, got \n%s\n", expectedStr, resultStr)
	}
}

func runGeneratePathTreeTest(t *testing.T, name string, paths []string, expectedFile string) {
	result := formatting.GeneratePathTree(paths)
	expected, err := os.ReadFile(filepath.Join("../test_data", expectedFile))
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}

	expectedStr := strings.TrimSpace(string(expected))
	resultStr := strings.TrimSpace(result)

	if resultStr != expectedStr {
		t.Errorf("%s: Expected:\n%s\n\nGot:\n%s", name, expectedStr, resultStr)
	}
}

func TestGeneratePathTreeBasicStructure(t *testing.T) {
	paths := []string{
		"cmd/ai-code-review/main.go",
		"go.mod",
		"internal/files/filtering.go",
		"internal/formatting/format.go",
	}
	runGeneratePathTreeTest(t, "Basic structure", paths, "basic_structure.txt")
}

func TestGeneratePathTreeEmptyDirectories(t *testing.T) {
	paths := []string{
		"dir1/file1.txt",
		"dir2/",
		"dir3/subdir/",
		"file2.txt",
	}
	runGeneratePathTreeTest(t, "Empty directories", paths, "empty_directories.txt")
}

func TestGeneratePathTreeMixedDepth(t *testing.T) {
	paths := []string{
		"a/very/deep/path/file.txt",
		"b/file.txt",
		"c/",
		"root_file.txt",
	}
	runGeneratePathTreeTest(t, "Mixed depth", paths, "mixed_depth.txt")
}

func TestGeneratePathTreeDuplicateParentDirectories(t *testing.T) {
	paths := []string{
		"parent/child1/file1.txt",
		"parent/child1/file2.txt",
		"parent/child2/file3.txt",
	}
	runGeneratePathTreeTest(t, "Duplicate parent directories", paths, "duplicate_parent_directories.txt")
}
