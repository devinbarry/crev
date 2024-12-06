package files_test

import (
	"github.com/devinbarry/crev/internal/files"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

// TestExplicitFilesPriority tests various scenarios where explicit files should
// override exclude patterns, including nested directories, multiple patterns,
// and mixed file types.
func TestExplicitFilesPriority(t *testing.T) {
	rootDir := t.TempDir()

	// Create a complex file structure
	fileStructure := map[string]string{
		"src/file1.go":             "content1",
		"src/file2.go":             "content2",
		"src/nested/file3.go":      "content3",
		"src/nested/file4.txt":     "content4",
		"build/output1.go":         "output1",
		"build/output2.txt":        "output2",
		"docs/readme.md":           "readme",
		"docs/api/overview.md":     "api docs",
		"vendor/lib1/module.go":    "module1",
		"vendor/lib2/package.json": "package",
	}
	createFiles(t, rootDir, fileStructure)

	tests := []struct {
		name            string
		excludePatterns []string
		explicitFiles   []string
		expectedFiles   []string
	}{
		{
			name:            "explicit files from excluded directory",
			excludePatterns: []string{"build/**"},
			explicitFiles: []string{
				filepath.Join(rootDir, "build/output1.go"),
				filepath.Join(rootDir, "build/output2.txt"),
			},
			expectedFiles: []string{
				filepath.Join(rootDir, "build/output1.go"),
				filepath.Join(rootDir, "build/output2.txt"),
				filepath.Join(rootDir, "docs"),
				filepath.Join(rootDir, "docs/api"),
				filepath.Join(rootDir, "docs/api/overview.md"),
				filepath.Join(rootDir, "docs/readme.md"),
				filepath.Join(rootDir, "src"),
				filepath.Join(rootDir, "src/file1.go"),
				filepath.Join(rootDir, "src/file2.go"),
				filepath.Join(rootDir, "src/nested"),
				filepath.Join(rootDir, "src/nested/file3.go"),
				filepath.Join(rootDir, "src/nested/file4.txt"),
				filepath.Join(rootDir, "vendor"),
				filepath.Join(rootDir, "vendor/lib1"),
				filepath.Join(rootDir, "vendor/lib1/module.go"),
				filepath.Join(rootDir, "vendor/lib2"),
				filepath.Join(rootDir, "vendor/lib2/package.json"),
			},
		},
		{
			name: "explicit files with multiple exclude patterns",
			excludePatterns: []string{
				"**/*.txt",
				"**/*.md",
				"vendor/**",
			},
			explicitFiles: []string{
				filepath.Join(rootDir, "src/nested/file4.txt"),
				filepath.Join(rootDir, "docs/readme.md"),
				filepath.Join(rootDir, "vendor/lib1/module.go"),
			},
			expectedFiles: []string{
				filepath.Join(rootDir, "build"),
				filepath.Join(rootDir, "build/output1.go"),
				filepath.Join(rootDir, "docs"),
				filepath.Join(rootDir, "docs/readme.md"),
				filepath.Join(rootDir, "src"),
				filepath.Join(rootDir, "src/file1.go"),
				filepath.Join(rootDir, "src/file2.go"),
				filepath.Join(rootDir, "src/nested"),
				filepath.Join(rootDir, "src/nested/file3.go"),
				filepath.Join(rootDir, "src/nested/file4.txt"),
				filepath.Join(rootDir, "vendor/lib1/module.go"),
			},
		},
		{
			name: "explicit files with extension and directory excludes",
			excludePatterns: []string{
				"**/*.go",
				"docs/**",
			},
			explicitFiles: []string{
				filepath.Join(rootDir, "src/file1.go"),
				filepath.Join(rootDir, "docs/api/overview.md"),
			},
			expectedFiles: []string{
				filepath.Join(rootDir, "build"),
				filepath.Join(rootDir, "build/output2.txt"),
				filepath.Join(rootDir, "docs/api/overview.md"),
				filepath.Join(rootDir, "src"),
				filepath.Join(rootDir, "src/file1.go"),
				filepath.Join(rootDir, "src/nested"),
				filepath.Join(rootDir, "src/nested/file4.txt"),
				filepath.Join(rootDir, "vendor"),
				filepath.Join(rootDir, "vendor/lib1"),
				filepath.Join(rootDir, "vendor/lib2"),
				filepath.Join(rootDir, "vendor/lib2/package.json"),
			},
		},
		{
			name:            "non-existent explicit files",
			excludePatterns: []string{"src/**"},
			explicitFiles: []string{
				filepath.Join(rootDir, "src/file1.go"),
				filepath.Join(rootDir, "non-existent.txt"),
			},
			expectedFiles: []string{
				filepath.Join(rootDir, "build"),
				filepath.Join(rootDir, "build/output1.go"),
				filepath.Join(rootDir, "build/output2.txt"),
				filepath.Join(rootDir, "docs"),
				filepath.Join(rootDir, "docs/api"),
				filepath.Join(rootDir, "docs/api/overview.md"),
				filepath.Join(rootDir, "docs/readme.md"),
				filepath.Join(rootDir, "src/file1.go"),
				filepath.Join(rootDir, "vendor"),
				filepath.Join(rootDir, "vendor/lib1"),
				filepath.Join(rootDir, "vendor/lib1/module.go"),
				filepath.Join(rootDir, "vendor/lib2"),
				filepath.Join(rootDir, "vendor/lib2/package.json"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePaths, err := files.GetAllFilePaths(rootDir, nil, tt.excludePatterns, tt.explicitFiles)
			require.NoError(t, err, "GetAllFilePaths failed")
			require.ElementsMatch(t, tt.expectedFiles, filePaths, "Incorrect paths returned")
		})
	}
}
