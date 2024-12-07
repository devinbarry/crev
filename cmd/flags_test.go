package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFileFlag tests the --file flag works correctly in various scenarios
func TestFileFlag(t *testing.T) {
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

	t.Run("single file", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--files", "main.go")
		require.NoError(t, err)

		expectedFiles := []string{"main.go"}
		unexpectedFiles := []string{"helper.go", "util/utils.go"}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("multiple files", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--files", "main.go", "--files", "helper.go")
		require.NoError(t, err)

		expectedFiles := []string{"main.go", "helper.go"}
		unexpectedFiles := []string{"util/utils.go"}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("file in subdirectory", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--files", "util/utils.go")
		require.NoError(t, err)

		expectedFiles := []string{"util/utils.go"}
		unexpectedFiles := []string{"main.go", "helper.go"}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("non-existent file", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--files", "nonexistent.go")
		require.Error(t, err, "Should fail when specified file doesn't exist")
		env.assertLogContains("no files found to bundle")
	})
}

// TestIncludeFlag tests the --include flag works correctly in various scenarios
func TestIncludeFlag(t *testing.T) {
	env := newTestEnv(t)

	// Setup a test project structure
	files := map[string]string{
		"src/main.go":        "package main",
		"src/util/helper.go": "package util",
		"test/main_test.go":  "package test",
		"docs/api.md":        "# API Docs",
		"docs/guide.md":      "# User Guide",
		"build/output.txt":   "build output",
	}
	env.createProjectStructure(files)

	t.Run("include single directory", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--include", "src/**")
		require.NoError(t, err)

		expectedFiles := []string{"src/main.go", "src/util/helper.go"}
		unexpectedFiles := []string{"test/main_test.go", "docs/api.md"}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("include by extension", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--include", "**/*.go")
		require.NoError(t, err)

		expectedFiles := []string{"src/main.go", "src/util/helper.go", "test/main_test.go"}
		unexpectedFiles := []string{"docs/api.md", "build/output.txt"}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("include multiple patterns", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--include", "src/**", "--include", "docs/**")
		require.NoError(t, err)

		expectedFiles := []string{
			"src/main.go",
			"src/util/helper.go",
			"docs/api.md",
			"docs/guide.md",
		}
		unexpectedFiles := []string{"test/main_test.go", "build/output.txt"}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("include pattern with no matches", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--include", "nonexistent/**")
		require.Error(t, err)
		env.assertLogContains("no files found to bundle")
	})
}

// TestExcludeFlag tests the --exclude flag works correctly in various scenarios
func TestExcludeFlag(t *testing.T) {
	env := newTestEnv(t)

	// Setup a test project structure
	files := map[string]string{
		"src/main.go":             "package main",
		"src/util/helper.go":      "package util",
		"src/util/helper_test.go": "package util_test",
		"test/main_test.go":       "package test",
		"docs/api.md":             "# API Docs",
		"build/output.txt":        "build output",
		".git/config":             "git config",
		"node_modules/pkg.js":     "module.exports = {};",
	}
	env.createProjectStructure(files)

	t.Run("exclude single directory", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--exclude", "test/**")
		require.NoError(t, err)

		expectedFiles := []string{"src/main.go", "docs/api.md"}
		unexpectedFiles := []string{"test/main_test.go"}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("exclude by extension", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--exclude", "**/*.md")
		require.NoError(t, err)

		expectedFiles := []string{"src/main.go", "test/main_test.go"}
		unexpectedFiles := []string{"docs/api.md"}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("exclude test files", func(t *testing.T) {
		err := env.executeBundleCmd(".", "--exclude", "**/*_test.go")
		require.NoError(t, err)

		expectedFiles := []string{"src/main.go", "src/util/helper.go"}
		unexpectedFiles := []string{
			"src/util/helper_test.go",
			"test/main_test.go",
		}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("exclude multiple patterns", func(t *testing.T) {
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
	})
}

// TestFlagCombinations tests various combinations of --file, --include, and --exclude flags
func TestFlagCombinations(t *testing.T) {
	env := newTestEnv(t)

	// Setup a test project structure
	files := map[string]string{
		"src/main.go":         "package main",
		"src/util/helper.go":  "package util",
		"test/main_test.go":   "package test",
		"test/helper_test.go": "package test",
		"docs/api.md":         "# API Docs",
		"build/output.txt":    "build output",
	}
	env.createProjectStructure(files)

	t.Run("include and exclude", func(t *testing.T) {
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
	})

	t.Run("explicit files override excludes", func(t *testing.T) {
		err := env.executeBundleCmd(".",
			"--files", "test/main_test.go",
			"--exclude", "**/*_test.go")
		require.NoError(t, err)

		expectedFiles := []string{"test/main_test.go"}
		unexpectedFiles := []string{
			"test/helper_test.go",
			"src/main.go",
		}
		env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)
	})

	t.Run("files with includes", func(t *testing.T) {
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
	})

	t.Run("complex pattern combination", func(t *testing.T) {
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
	})
}
