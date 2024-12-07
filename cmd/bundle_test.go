package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBundleCommandBasic tests the basic functionality of the bundle command without any exclude patterns.
func TestBundleCommandBasic(t *testing.T) {
	// Create a mock project structure
	/*
		/tempDir
		├── include.go
		├── main.go
		└── internal
		    ├── files
		    │   └── reading.go
		    └── formatting
		        └── format.go
	*/
	env := newTestEnv(t)
	files := map[string]string{
		"include.go":                    "package main",
		"main.go":                       "package main",
		"internal/files/reading.go":     "package files",
		"internal/formatting/format.go": "package formatting",
	}
	env.createProjectStructure(files)

	// Setup basic config
	env.setupConfig(basicConfig)

	// Execute bundle command
	err := env.executeBundleCmd(".")
	require.NoError(t, err, "Bundle command execution failed")

	// Verify output file contents
	expectedFiles := []string{
		"include.go",
		"main.go",
		"internal/files/reading.go",
		"internal/formatting/format.go",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, nil)

	// Verify log messages
	env.assertLogContains("Project overview successfully saved to:")
}

// TestBundleCommandWithConfigExcludes tests the bundle command with exclude patterns from config.
func TestBundleCommandWithConfigExcludes(t *testing.T) {
	// Create a mock project structure
	/*
		/tempDir
		├── include.go
		├── exclude.md
		├── build_something.py
		├── .git/
		│   └── config
		├── node_modules/
		│   └── module.js
		└── images/
		    └── image.png
	*/
	env := newTestEnv(t)
	files := map[string]string{
		"include.go":             "package main",
		"exclude.md":             "# Exclude",
		"build_something.py":     "# Python build script",
		".git/config":            "[core]",
		"node_modules/module.js": "// JS Module",
		"images/image.png":       "PNGDATA",
	}
	env.createProjectStructure(files)

	// Setup full config file with all exclusions
	env.setupConfig(fullConfig)

	// Execute bundle command
	err := env.executeBundleCmd(".")
	require.NoError(t, err, "Bundle command execution failed")

	// Verify included and excluded files
	expectedFiles := []string{
		"include.go",
		"build_something.py",
	}
	unexpectedFiles := []string{
		"images/image.png",
		"exclude.md",
		".git/config",
		"node_modules/module.js",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)

	// Verify log messages
	env.assertLogContains("Project overview successfully saved to:")
}

// TestBundleCommandWithExplicitFiles tests that explicitly included files are included even if they match exclude patterns.
func TestBundleCommandWithExplicitFiles(t *testing.T) {
	env := newTestEnv(t)

	// Create files
	files := map[string]string{
		"include.go":         "package main",
		"exclude.md":         "# Exclude",
		"build_something.py": "# Python build script",
		"images/image.png":   "PNGDATA",
	}
	env.createProjectStructure(files)

	// Setup config excluding .md files
	env.setupConfig(excludeConfig)

	// Execute bundle command with explicit file inclusion
	err := env.executeBundleCmd(".", "--files", "exclude.md")
	require.NoError(t, err, "Bundle command execution failed")

	// Verify file contents
	expectedFiles := []string{
		"include.go",
		"build_something.py",
		"exclude.md", // Should be included due to explicit inclusion
	}
	unexpectedFiles := []string{
		"images/image.png", // Images are always excluded because they are not text
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)

	// Verify log messages
	env.assertLogContains("Project overview successfully saved to:")
}

// New test to ensure that explicitly included files always override exclude patterns
func TestBundleCommandWithExplicitFilesOverridesExclude(t *testing.T) {
	env := newTestEnv(t)

	// Create test files
	files := map[string]string{
		"pycharm.logs.bak": "logger.info",
	}
	env.createProjectStructure(files)

	// Setup config excluding .bak files
	configContent := `
include:
  - "**/*"
exclude:
  - "*.bak"
`
	env.setupConfig(configContent)

	// Execute bundle command with explicit file inclusion
	err := env.executeBundleCmd(".", "--files", "pycharm.logs.bak")
	require.NoError(t, err, "Bundle command execution failed")

	// Verify explicitly included file is present
	expectedFiles := []string{
		"pycharm.logs.bak",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, nil)
}

// TestBundleCommandWithNoFiles tests the bundle command when no files are included due to exclude patterns.
func TestBundleCommandWithNoFiles(t *testing.T) {
	// Create a mock project structure
	/*
		/tempDir
		├── exclude.md
		└── .git/
		    └── config
	*/
	env := newTestEnv(t)

	// Create test files that will all be excluded
	files := map[string]string{
		"exclude.md":  "# Exclude",
		".git/config": "[core]",
	}
	env.createProjectStructure(files)

	// Setup config that excludes all files
	configContent := `
include: []
exclude:
  - "**/*"
`
	env.setupConfig(configContent)

	// Execute bundle command
	err := env.executeBundleCmd(".")
	env.assertErrorContains(err, "no files found to bundle")
}

// TestBundleCommandWithIncludeAndExcludePatterns tests combining include and exclude patterns.
func TestBundleCommandWithIncludeAndExcludePatterns(t *testing.T) {
	// Create a mock project structure
	/*
		/tempDir
		├── src/
		│   ├── main.go
		│   └── utils.go
		├── vendor/
		│   └── lib.go
		├── test/
		│   └── main_test.go
		└── README.md
	*/
	env := newTestEnv(t)

	// Create test files
	files := map[string]string{
		"src/main.go":       "package main",
		"src/utils.go":      "package utils",
		"vendor/lib.go":     "package lib",
		"test/main_test.go": "package main_test",
		"README.md":         "# Project",
	}
	env.createProjectStructure(files)

	// Setup config with include and exclude patterns
	configContent := `
include:
  - "src/**"
  - "test/**"
exclude:
  - "vendor/**"
  - "*.md"
`
	env.setupConfig(configContent)

	// Execute bundle command
	err := env.executeBundleCmd(".")
	require.NoError(t, err, "Bundle command execution failed")

	// Verify included and excluded files
	expectedFiles := []string{
		"src/main.go",
		"src/utils.go",
		"test/main_test.go",
	}
	unexpectedFiles := []string{
		"vendor/lib.go",
		"README.md",
	}
	env.assertFileContents("crev-project.txt", expectedFiles, unexpectedFiles)

	// Verify log messages
	env.assertLogContains("Project overview successfully saved to:")
}

// TestBundleCommandHandlesNonExistentPath tests the bundle command when the specified path does not exist.
func TestBundleCommandHandlesNonExistentPath(t *testing.T) {
	env := newTestEnv(t)

	// Setup config (even though we won't have any files)
	env.setupConfig(basicConfig)

	// Use a non-existent directory path
	nonExistentDir := filepath.Join(os.TempDir(), "non_existent_dir_123456")

	// Execute bundle command with non-existent path
	err := env.executeBundleCmd(nonExistentDir)
	env.assertErrorContains(err, "does not exist")
}
