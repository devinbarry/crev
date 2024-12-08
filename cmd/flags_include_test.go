package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSingleFileFlag tests bundling a single file
func TestSingleFileFlag(t *testing.T) {
	env := newTestEnv(t)

	// Setup a test project structure
	files := map[string]string{
		"main.go":        "package main",
		"helper.go":      "package main",
		"util/utils.go":  "package util",
		"docs/readme.md": "# Documentation",
		"test/test.go":   "package test",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--files", "main.go")
	require.NoError(t, err)

	expectedFiles := []string{"main.go", "helper.go", "util/utils.go"}
	env.assertFileContents("crev-project.txt", expectedFiles, nil)
}

// TestMultipleFilesFlag tests bundling multiple specified files
func TestMultipleFilesFlag(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"main.go":        "package main",
		"helper.go":      "package main",
		"util/utils.go":  "package util",
		"docs/readme.md": "# Documentation",
		"test/test.go":   "package test",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--files", "main.go", "--files", "helper.go")
	require.NoError(t, err)

	expectedFiles := []string{"main.go", "helper.go", "util/utils.go"}
	env.assertFileContents("crev-project.txt", expectedFiles, nil)
}

// TestFileInSubdirectoryFlag tests bundling a file from a subdirectory
func TestFileInSubdirectoryFlag(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"main.go":        "package main",
		"helper.go":      "package main",
		"util/utils.go":  "package util",
		"docs/readme.md": "# Documentation",
		"test/test.go":   "package test",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--files", "util/utils.go")
	require.NoError(t, err)

	expectedFiles := []string{"util/utils.go", "main.go", "helper.go"}
	env.assertFileContents("crev-project.txt", expectedFiles, nil)
}

// TestNonExistentFileFlag tests behavior with non-existent file
func TestNonExistentFileFlag(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"main.go": "package main",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--files", "nonexistent.go")
	// TODO We don't check if files exist yet
	require.Error(t, err, "Should fail when specified file doesn't exist")
	env.assertLogContains("no files found to bundle")
}

// TestIncludeSingleDirectory tests including files from a single directory
func TestIncludeSingleDirectory(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":        "package main",
		"src/util/helper.go": "package util",
		"test/main_test.go":  "package test",
		"docs/api.md":        "# API Docs",
		"docs/guide.md":      "# User Guide",
		"build/output.txt":   "build output",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--include", "src/**")
	require.NoError(t, err)

	expectedFiles := []string{"src/main.go", "src/util/helper.go"}
	unexpectedFiles := []string{"test/main_test.go", "docs/api.md"}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestIncludeByExtension tests including files by extension pattern
func TestIncludeByExtension(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":        "package main",
		"src/util/helper.go": "package util",
		"test/main_test.go":  "package test",
		"docs/api.md":        "# API Docs",
		"build/output.txt":   "build output",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--include", "**/*.go")
	require.NoError(t, err)

	expectedFiles := []string{"src/main.go", "src/util/helper.go", "test/main_test.go"}
	unexpectedFiles := []string{"docs/api.md", "build/output.txt"}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestIncludeMultiplePatterns tests including files using multiple patterns
func TestIncludeMultiplePatterns(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":        "package main",
		"src/util/helper.go": "package util",
		"test/main_test.go":  "package test",
		"docs/api.md":        "# API Docs",
		"docs/guide.md":      "# User Guide",
		"build/output.txt":   "build output",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--include", "src/**", "--include", "docs/**")
	require.NoError(t, err)

	expectedFiles := []string{
		"src/main.go",
		"src/util/helper.go",
		"docs/api.md",
		"docs/guide.md",
	}
	unexpectedFiles := []string{
		"test/main_test.go",
		"test/main_test.go",
		"build/output.txt",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestIncludeNoMatches tests behavior when include pattern matches no files
func TestIncludeNoMatches(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go": "package main",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--include", "nonexistent/**")
	// No files are included in this bundle so an error is raised saying so.
	env.assertErrorContains(err, "no files found to bundle.")
}
