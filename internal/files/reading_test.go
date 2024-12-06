package files

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

// Test the preprocessExcludePatterns function
func TestPreprocessExcludePatterns(t *testing.T) {
	rootDir := t.TempDir()

	// Create test filesystem structure
	require.NoError(t, os.Mkdir(filepath.Join(rootDir, "dir"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(rootDir, "file.txt"), []byte("content"), 0644))

	testCases := []struct {
		name     string
		pattern  string
		expected []string // Can match multiple patterns
	}{
		{
			name:     "existing directory with trailing slash",
			pattern:  "dir/",
			expected: []string{"dir", "dir/**"}, // Should match both dir and contents
		},
		{
			name:     "existing directory without trailing slash",
			pattern:  "dir",
			expected: []string{"dir", "dir/**"}, // Should behave same as with slash
		},
		{
			name:     "existing file",
			pattern:  "file.txt",
			expected: []string{"file.txt"}, // Should remain unchanged
		},
		{
			name:     "non-existent directory pattern",
			pattern:  "nonexistent/",
			expected: []string{"nonexistent"}, // Should clean trailing slash
		},
		{
			name:     "simple string pattern",
			pattern:  "*.md",
			expected: []string{"*.md"}, // Should remain unchanged
		},
		{
			name:     "empty pattern",
			pattern:  "",
			expected: []string{}, // Should be skipped
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := preprocessExcludePatterns(rootDir, []string{tc.pattern})

			// Check that all expected patterns are present
			for _, exp := range tc.expected {
				found := false
				for _, res := range result {
					if res == exp {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("pattern %q not found in result %v", exp, result)
				}
			}

			// Check no unexpected patterns are present
			if len(result) != len(tc.expected) {
				t.Errorf("got unexpected number of patterns, expected %v, got %v", tc.expected, result)
			}
		})
	}
}
