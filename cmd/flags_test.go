package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestExcludeSingleDirectory tests excluding a single directory
func TestExcludeSingleDirectory(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":             "package main",
		"src/util/helper.go":      "package util",
		"src/util/helper_test.go": "package util_test",
		"test/main_test.go":       "package test",
		"docs/api.md":             "# API Docs",
		"build/output.txt":        "build output",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--exclude", "test/**")
	require.NoError(t, err)

	expectedFiles := []string{
		"src/main.go",
		"src/util/helper.go",
		"src/util/helper_test.go",
		"docs/api.md",
		"build/output.txt",
	}
	unexpectedFiles := []string{"test/main_test.go"}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestExcludeByExtension tests excluding files by extension
func TestExcludeByExtension(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":        "package main",
		"src/util/helper.go": "package util",
		"test/main_test.go":  "package test",
		"docs/api.md":        "# API Docs",
		"build/output.txt":   "build output",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--exclude", "**/*.md")
	require.NoError(t, err)

	expectedFiles := []string{"src/main.go", "test/main_test.go"}
	unexpectedFiles := []string{"docs/api.md"}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestExcludeTestFiles tests excluding test files
func TestExcludeTestFiles(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":             "package main",
		"src/util/helper.go":      "package util",
		"src/util/helper_test.go": "package util_test",
		"test/main_test.go":       "package test",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--exclude", "**/*_test.go")
	require.NoError(t, err)

	expectedFiles := []string{"src/main.go", "src/util/helper.go"}
	unexpectedFiles := []string{
		"src/util/helper_test.go",
		"test/main_test.go",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestExcludeMultiplePatterns tests excluding files using multiple patterns
func TestExcludeMultiplePatterns(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":        "package main",
		"src/util/helper.go": "package util",
		"test/main_test.go":  "package test",
		"docs/api.md":        "# API Docs",
		"build/output.txt":   "build output",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".",
		"--exclude", "test/**",
		"--exclude", "docs/**",
		"--exclude", "build/**")
	require.NoError(t, err)

	expectedFiles := []string{"src/main.go", "src/util/helper.go"}
	unexpectedFiles := []string{
		"test/main_test.go",
		"docs/api.md",
		"build/output.txt",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestIncludeAndExclude tests combining include and exclude patterns
func TestIncludeAndExclude(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":         "package main",
		"src/util/helper.go":  "package util",
		"test/main_test.go":   "package test",
		"test/helper_test.go": "package test",
		"docs/api.md":         "# API Docs",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".",
		"--include", "src/**",
		"--include", "test/**",
		"--exclude", "**/*_test.go")
	require.NoError(t, err)

	expectedFiles := []string{"src/main.go", "src/util/helper.go"}
	unexpectedFiles := []string{
		"test/main_test.go",
		"test/helper_test.go",
		"docs/api.md",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestFilesOverrideExcludes tests that explicit files override exclude patterns
func TestFilesOverrideExcludes(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/bundle.go":       "package main",
		"src/bundle_test.go":  "package main",
		"src/main.go":         "package main",
		"test/main_test.go":   "package test",
		"test/helper_test.go": "package test",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".", "--files", "test/main_test.go", "--exclude", "**/*_test.go")
	require.NoError(t, err)

	expectedFiles := []string{"test/main_test.go"}
	unexpectedFiles := []string{
		"src/bundle.go",
		"src/bundle_test.go",
		"src/main.go",
		"test/helper_test.go",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestFilesWithIncludes tests combining explicit files with include patterns
func TestFilesWithIncludes(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":        "package main",
		"src/util/helper.go": "package util",
		"test/main_test.go":  "package test",
		"docs/api.md":        "# API Docs",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".",
		"--files", "docs/api.md",
		"--include", "src/**")
	require.NoError(t, err)

	expectedFiles := []string{
		"docs/api.md",
		"src/main.go",
		"src/util/helper.go",
	}
	unexpectedFiles := []string{"test/main_test.go"}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}

// TestComplexPatternCombination tests complex combination of files, includes, and excludes
func TestComplexPatternCombination(t *testing.T) {
	env := newTestEnv(t)
	files := map[string]string{
		"src/main.go":         "package main",
		"src/util/helper.go":  "package util",
		"test/main_test.go":   "package test",
		"test/helper_test.go": "package test",
		"docs/api.md":         "# API Docs",
	}
	env.createProjectStructure(files)

	err := env.executeBundleCmd(".",
		"--files", "docs/api.md",
		"--include", "src/**",
		"--include", "test/**",
		"--exclude", "**/*_test.go",
		"--exclude", "**/helper.go")
	require.NoError(t, err)

	expectedFiles := []string{"docs/api.md", "src/main.go"}
	unexpectedFiles := []string{
		"src/util/helper.go",
		"test/main_test.go",
		"test/helper_test.go",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
}
