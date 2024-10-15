package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPreprocessExcludePatterns(t *testing.T) {
	rootDir := t.TempDir()

	// Create a directory and a file
	os.Mkdir(filepath.Join(rootDir, "dir"), 0755)
	os.WriteFile(filepath.Join(rootDir, "file.txt"), []byte("content"), 0644)

	// Prepare exclude patterns
	excludePatterns := []string{"dir/", "file.txt", "nonexistent/", "empty_string", ""}
	expectedPatterns := []string{"dir/**", "file.txt", "nonexistent/", "empty_string", ""}

	processedPatterns := preprocessExcludePatterns(rootDir, excludePatterns)

	if len(processedPatterns) != len(expectedPatterns) {
		t.Fatalf("expected %d patterns, got %d", len(expectedPatterns), len(processedPatterns))
	}

	for i, exp := range expectedPatterns {
		if processedPatterns[i] != exp {
			t.Errorf("expected pattern %q, got %q", exp, processedPatterns[i])
		}
	}
}
