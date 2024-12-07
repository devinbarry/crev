package files_test

import (
	"github.com/devinbarry/crev/internal/files"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

// TestExplicitFilesPriority_ExplicitFromExcludedDirectory tests that explicitly specified files
// are included in the results even if they are in directories that would otherwise be excluded.
// This ensures that the --files flag takes precedence over exclude patterns.
func TestExplicitFilesPriority_ExplicitFromExcludedDirectory(t *testing.T) {
	// t.TempDir() is a Go testing helper that creates a temporary directory
	// and automatically cleans it up after the test completes
	rootDir := t.TempDir()

	// Create a test file structure that mimics a typical project layout
	// with source files, build outputs, documentation, and vendor directories
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

	// Exclude all files in the build directory
	excludePatterns := []string{"build/**"}

	// Explicitly include specific files from the excluded build directory
	explicitFiles := []string{
		filepath.Join(rootDir, "build/output1.go"),
		filepath.Join(rootDir, "build/output2.txt"),
	}

	// Define the complete set of files we expect to see in the results
	// This includes:
	// 1. Explicitly included files (even from excluded directories)
	// 2. All directories (they're needed for traversal)
	// 3. All files not in excluded directories
	expectedFiles := []string{
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
	}

	// Get all file paths using the function under test
	// nil is passed as includePatterns, meaning all files are included by default
	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, explicitFiles)
	require.NoError(t, err, "GetAllFilePaths failed")

	// require.ElementsMatch checks that two slices contain the same elements, regardless of order
	require.ElementsMatch(t, expectedFiles, filePaths, "Incorrect paths returned")

	// Verify that non-existent files are not included in the results
	require.NotContains(t, filePaths, filepath.Join(rootDir, "non-existent-file.txt"))
}

// TestExplicitFilesPriority_MultipleExcludePatterns tests that multiple exclude patterns
// work together correctly while still respecting explicitly included files.
func TestExplicitFilesPriority_MultipleExcludePatterns(t *testing.T) {
	rootDir := t.TempDir()

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

	// Test multiple exclude patterns:
	// 1. All .txt files
	// 2. All .md files
	// 3. All files in vendor directory
	excludePatterns := []string{
		"**/*.txt",  // ** matches any number of directories
		"**/*.md",   // * matches any characters except path separator
		"vendor/**", // Exclude all vendor directory contents
	}

	// Explicitly include some files that would otherwise be excluded
	explicitFiles := []string{
		filepath.Join(rootDir, "src/nested/file4.txt"),
		filepath.Join(rootDir, "docs/readme.md"),
		filepath.Join(rootDir, "vendor/lib1/module.go"),
	}

	// Expected files include:
	// 1. All .go files (except in vendor)
	// 2. Explicitly included files
	// 3. Necessary directory structure
	expectedFiles := []string{
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
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, explicitFiles)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expectedFiles, filePaths, "Incorrect paths returned")

	// Verify that excluded files that weren't explicitly included remain excluded
	require.NotContains(t, filePaths, filepath.Join(rootDir, "docs/api/overview.md"))
	require.NotContains(t, filePaths, filepath.Join(rootDir, "vendor/lib2/package.json"))
}

// TestExplicitFilesPriority_ExtensionAndDirectoryExcludes tests the interaction between
// file extension patterns and directory patterns, along with explicit file overrides.
func TestExplicitFilesPriority_ExtensionAndDirectoryExcludes(t *testing.T) {
	rootDir := t.TempDir()

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

	// Test combination of extension and directory exclusions
	excludePatterns := []string{
		"**/*.go", // Exclude all .go files
		"docs/**", // Exclude all files in docs directory
	}

	// Explicitly include some excluded files
	explicitFiles := []string{
		filepath.Join(rootDir, "src/file1.go"),
		filepath.Join(rootDir, "docs/api/overview.md"),
	}

	// Expected files include:
	// 1. Non-excluded files (.txt and .json)
	// 2. Explicitly included files
	// 3. Necessary directory structure
	expectedFiles := []string{
		filepath.Join(rootDir, "build"),
		filepath.Join(rootDir, "build/output2.txt"),
		filepath.Join(rootDir, "docs/api/overview.md"),
		filepath.Join(rootDir, "src"),
		filepath.Join(rootDir, "src/file1.go"),
		filepath.Join(rootDir, "src/nested"),
		filepath.Join(rootDir, "src/nested/file4.txt"),
		filepath.Join(rootDir, "vendor"),
		filepath.Join(rootDir, "vendor/lib2"),
		filepath.Join(rootDir, "vendor/lib2/package.json"),
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, explicitFiles)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expectedFiles, filePaths, "Incorrect paths returned")

	// Verify that excluded files that weren't explicitly included remain excluded
	require.NotContains(t, filePaths, filepath.Join(rootDir, "src/file2.go"))
	require.NotContains(t, filePaths, filepath.Join(rootDir, "docs/readme.md"))
}

// TestExplicitFilesPriority_NonExistentExplicitFiles tests that the system handles
// non-existent explicit files gracefully while still processing valid files correctly.
func TestExplicitFilesPriority_NonExistentExplicitFiles(t *testing.T) {
	rootDir := t.TempDir()

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

	// Exclude all files in src directory
	excludePatterns := []string{"src/**"}

	// Include one valid file and one non-existent file
	explicitFiles := []string{
		filepath.Join(rootDir, "src/file1.go"),
		filepath.Join(rootDir, "non-existent.txt"), // This file doesn't exist
	}

	// Expected files include:
	// 1. All non-excluded files
	// 2. Valid explicitly included files (but not non-existent ones)
	// 3. Necessary directory structure
	expectedFiles := []string{
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
	}

	filePaths, err := files.GetAllFilePaths(rootDir, nil, excludePatterns, explicitFiles)
	require.NoError(t, err, "GetAllFilePaths failed")
	require.ElementsMatch(t, expectedFiles, filePaths, "Incorrect paths returned")

	// Verify that excluded files and non-existent files are not included
	require.NotContains(t, filePaths, filepath.Join(rootDir, "src/file2.go"))
	require.NotContains(t, filePaths, filepath.Join(rootDir, "non-existent.txt"))
}
