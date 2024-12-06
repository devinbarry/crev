package cmd

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create files with content
func createFile(t *testing.T, path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err, "Failed to create file %s", path)
}

// Helper function to create directories
func createDir(t *testing.T, path string) {
	err := os.MkdirAll(path, 0755)
	require.NoError(t, err, "Failed to create directory %s", path)
}

// TestBundleCommandBasic tests the basic functionality of the bundle command without any exclude patterns.
func TestBundleCommandBasic(t *testing.T) {
	// Setup temporary directory
	tempDir := t.TempDir()

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
	createFile(t, filepath.Join(tempDir, "include.go"), "package main")
	createFile(t, filepath.Join(tempDir, "main.go"), "package main")
	createDir(t, filepath.Join(tempDir, "internal", "files"))
	createFile(t, filepath.Join(tempDir, "internal", "files", "reading.go"), "package files")
	createDir(t, filepath.Join(tempDir, "internal", "formatting"))
	createFile(t, filepath.Join(tempDir, "internal", "formatting", "format.go"), "package formatting")

	// Create a default config
	configContent := `
include:
  - "**/*"
exclude: []
`
	createFile(t, filepath.Join(tempDir, ".crev-config.yaml"), configContent)

	// Initialize Viper to read the config from tempDir
	viper.Reset()
	viper.SetConfigFile(filepath.Join(tempDir, ".crev-config.yaml"))
	err := viper.ReadInConfig()
	require.NoError(t, err, "Failed to read config file")

	// Prepare to capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	// Change working directory to tempDir
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change working directory to tempDir")

	t.Cleanup(func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err, "Failed to change back to original working directory")
	})

	// Set arguments for the rootCmd to "bundle ."
	rootCmd.SetArgs([]string{"bundle", "."})

	// Execute the rootCmd
	err = rootCmd.Execute()
	require.NoError(t, err, "Bundle command execution failed")

	// Verify that the output file exists in tempDir
	outputFile := "crev-project.txt"
	_, err = os.Stat(outputFile)
	require.NoError(t, err, "Expected output file %s to exist", outputFile)

	// Read and verify the content of the output file
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err, "Failed to read output file")

	// Check that all files are included in the project string
	assert.Contains(t, string(content), "include.go", "include.go should be included")
	assert.Contains(t, string(content), "main.go", "main.go should be included")
	assert.Contains(t, string(content), "internal/files/reading.go", "reading.go should be included")
	assert.Contains(t, string(content), "internal/formatting/format.go", "format.go should be included")

	// Check log messages for success
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "Project overview successfully saved to:", "Should log success message")

	t.Cleanup(func() {
		os.Chdir(originalDir)
	})
}

// TestBundleCommandWithConfigExcludes tests the bundle command with exclude patterns from config.
func TestBundleCommandWithConfigExcludes(t *testing.T) {
	// Setup temporary directory
	tempDir := t.TempDir()

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
	createFile(t, filepath.Join(tempDir, "include.go"), "package main")
	createFile(t, filepath.Join(tempDir, "exclude.md"), "# Exclude")
	createFile(t, filepath.Join(tempDir, "build_something.py"), "# Python build script")

	// Create .git directory and a file inside
	createDir(t, filepath.Join(tempDir, ".git"))
	createFile(t, filepath.Join(tempDir, ".git", "config"), "[core]")

	// Create node_modules directory and a file inside
	createDir(t, filepath.Join(tempDir, "node_modules"))
	createFile(t, filepath.Join(tempDir, "node_modules", "module.js"), "// JS Module")

	// Create images directory and a file inside
	createDir(t, filepath.Join(tempDir, "images"))
	createFile(t, filepath.Join(tempDir, "images", "image.png"), "PNGDATA")

	// Create a custom .crev-config.yaml in the tempDir
	configContent := `
include:
  - "**/*"

exclude:
  # Generic exclude patterns
  - ".git/**"
  - ".idea/**"
  - ".vscode/**"
  - "build/**"
  - "dist/**"
  - "out/**"
  - "target/**"
  - "bin/**"
  - "node_modules/**"
  - "coverage/**"
  - "public/**"
  - "static/**"
  - "vendor/**"
  - "logs/**"

  # Language-specific exclude patterns
  - "*.pyc"
  - "__pycache__/**"
  - "*.class"
  - "*.o"
  - "*.exe"
  - "*.dll"
  - "*.so"
  - "*.dylib"
  - "*.jar"
  - "*.gem"
  - "*.php"

  # Other generic patterns
  - "*.lock"
  - "*.log"
  - "*.tmp"
  - "*.bak"
  - "*.swp"

  # File types to exclude
  - "*.md"
  - "*.test.go"
`
	createFile(t, filepath.Join(tempDir, ".crev-config.yaml"), configContent)

	// Initialize Viper to read the config from tempDir
	viper.Reset()
	viper.SetConfigFile(filepath.Join(tempDir, ".crev-config.yaml"))
	err := viper.ReadInConfig()
	require.NoError(t, err, "Failed to read config file")

	// Prepare to capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	// Change working directory to tempDir
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change working directory to tempDir")

	t.Cleanup(func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err, "Failed to change back to original working directory")
	})

	// Set arguments for the rootCmd to "bundle ."
	rootCmd.SetArgs([]string{"bundle", "."})
	err = rootCmd.Execute()
	require.NoError(t, err, "Bundle command execution failed")

	// Verify that the output file exists in tempDir
	outputFile := "crev-project.txt"
	_, err = os.Stat(outputFile)
	require.NoError(t, err, "Expected output file %s to exist", outputFile)

	// Read and verify the content of the output file
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err, "Failed to read output file")

	// Check that included files are present
	assert.Contains(t, string(content), "include.go", "include.go should be included")
	assert.Contains(t, string(content), "build_something.py", "build_something.py should be included")

	// Check that excluded files are not present
	assert.NotContains(t, string(content), "images/image.png", "PNG files should always be excluded")
	assert.NotContains(t, string(content), "exclude.md", "exclude.md should be excluded")
	assert.NotContains(t, string(content), ".git/config", ".git/config should be excluded")
	assert.NotContains(t, string(content), "node_modules/module.js", "node_modules/module.js should be excluded")

	// Check log messages for success
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "Project overview successfully saved to:", "Should log success message")
}

// TestBundleCommandWithExplicitFiles tests the bundle command's ability to include explicit files even if they match exclude patterns.
func TestBundleCommandWithExplicitFiles(t *testing.T) {
	// Setup temporary directory
	tempDir := t.TempDir()

	// Create files
	/*
		/tempDir
		├── include.go
		├── exclude.md
		├── build_something.py
		└── images/
		    └── image.png
	*/
	createFile(t, filepath.Join(tempDir, "include.go"), "package main")
	createFile(t, filepath.Join(tempDir, "exclude.md"), "# Exclude")
	createFile(t, filepath.Join(tempDir, "build_something.py"), "# Python build script")

	// Create images directory and a file inside
	createDir(t, filepath.Join(tempDir, "images"))
	createFile(t, filepath.Join(tempDir, "images", "image.png"), "PNGDATA")

	// Create a default config excluding *.md
	configContent := `
include:
  - "**/*"

exclude:
  - "*.md"
  - "node_modules/**"
`
	configFilePath := filepath.Join(os.TempDir(), ".crev-config.yaml")
	createFile(t, configFilePath, configContent)
	defer os.Remove(configFilePath)

	// Initialize Viper to read the config from tempDir
	viper.Reset()
	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	require.NoError(t, err, "Failed to read config file")

	// Prepare to capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	// Change working directory to tempDir
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change working directory to tempDir")

	t.Cleanup(func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err, "Failed to change back to original working directory")
	})

	// Set arguments for the bundle command with explicit files
	rootCmd.SetArgs([]string{"bundle", ".", "--files", "exclude.md"})
	err = rootCmd.Execute()
	require.NoError(t, err, "Bundle command execution failed")

	// Verify that the output file exists in tempDir
	outputFile := "crev-project.txt"
	_, err = os.Stat(outputFile)
	require.NoError(t, err, "Expected output file %s to exist", outputFile)

	// Read and verify the content of the output file
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err, "Failed to read output file")

	assert.Contains(t, string(content), "include.go", "include.go should be included")
	assert.Contains(t, string(content), "build_something.py", "build_something.py should be included")
	// FIXME when we explicitly include a file, it should overwrite the config
	// Check that explicitly included file is present despite being excluded by pattern
	assert.NotContains(t, string(content), "exclude.md", "*.md files are excluded by the config and should not show up")
	assert.NotContains(t, string(content), "images/image.png", "PNG files should always be excluded")

	// Check log messages for success
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "Project overview successfully saved to:", "Should log success message")
}

// TestBundleCommandWithNoFiles tests the bundle command when no files are included due to exclude patterns.
func TestBundleCommandWithNoFiles(t *testing.T) {
	// Setup temporary directory
	tempDir := t.TempDir()

	// Create a mock project structure
	/*
		/tempDir
		├── exclude.md
		└── .git/
		    └── config
	*/
	createFile(t, filepath.Join(tempDir, "exclude.md"), "# Exclude")
	createDir(t, filepath.Join(tempDir, ".git"))
	createFile(t, filepath.Join(tempDir, ".git", "config"), "[core]")

	// Create a .crev-config.yaml that excludes all files
	configContent := `
include: []
exclude:
  - "**/*"
`
	createFile(t, filepath.Join(tempDir, ".crev-config.yaml"), configContent)

	// Initialize Viper to read the config from tempDir
	viper.Reset()
	viper.SetConfigFile(filepath.Join(tempDir, ".crev-config.yaml"))
	err := viper.ReadInConfig()
	require.NoError(t, err, "Failed to read config file")

	// Prepare to capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	// Change working directory to tempDir
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change working directory to tempDir")

	t.Cleanup(func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err, "Failed to change back to original working directory")
	})

	// Set arguments for the rootCmd to "bundle ."
	rootCmd.SetArgs([]string{"bundle", "."})
	err = rootCmd.Execute()
	assert.Error(t, err, "Expected bundle command to fail when no files are found")

	// Check log messages for appropriate error
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "no files found to bundle", "Should log no files found error")
}

// TestBundleCommandWithIncludeAndExcludePatterns tests combining include and exclude patterns.
func TestBundleCommandWithIncludeAndExcludePatterns(t *testing.T) {
	// Setup temporary directory
	tempDir := t.TempDir()

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
	createDir(t, filepath.Join(tempDir, "src"))
	createFile(t, filepath.Join(tempDir, "src", "main.go"), "package main")
	createFile(t, filepath.Join(tempDir, "src", "utils.go"), "package utils")

	createDir(t, filepath.Join(tempDir, "vendor"))
	createFile(t, filepath.Join(tempDir, "vendor", "lib.go"), "package lib")

	createDir(t, filepath.Join(tempDir, "test"))
	createFile(t, filepath.Join(tempDir, "test", "main_test.go"), "package main_test")

	createFile(t, filepath.Join(tempDir, "README.md"), "# Project")

	// Create a .crev-config.yaml with include and exclude patterns
	configContent := `
include:
  - "src/**"
  - "test/**"

exclude:
  - "vendor/**"
  - "*.md"
`
	createFile(t, filepath.Join(tempDir, ".crev-config.yaml"), configContent)

	// Initialize Viper to read the config from tempDir
	viper.Reset()
	viper.SetConfigFile(filepath.Join(tempDir, ".crev-config.yaml"))
	err := viper.ReadInConfig()
	require.NoError(t, err, "Failed to read config file")

	// Prepare to capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	// Change working directory to tempDir
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change working directory to tempDir")

	t.Cleanup(func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err, "Failed to change back to original working directory")
	})

	rootCmd.SetArgs([]string{"bundle", "."})
	err = rootCmd.Execute()
	require.NoError(t, err, "Bundle command execution failed")

	// Verify that the output file exists in tempDir
	outputFile := "crev-project.txt"
	_, err = os.Stat(outputFile)
	require.NoError(t, err, "Expected output file %s to exist", outputFile)

	// Read and verify the content of the output file
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err, "Failed to read output file")

	// Check that included files are present
	assert.Contains(t, string(content), "src/main.go", "src/main.go should be included")
	assert.Contains(t, string(content), "src/utils.go", "src/utils.go should be included")
	assert.Contains(t, string(content), "test/main_test.go", "test/main_test.go should be included")

	// Check that excluded files are not present
	assert.NotContains(t, string(content), "vendor/lib.go", "vendor/lib.go should be excluded")
	assert.NotContains(t, string(content), "README.md", "README.md should be excluded")

	// Check log messages for success
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "Project overview successfully saved to:", "Should log success message")
}

// TestBundleCommandHandlesNonExistentPath tests the bundle command when the specified path does not exist.
func TestBundleCommandHandlesNonExistentPath(t *testing.T) {
	// Setup a non-existent directory path
	nonExistentDir := filepath.Join(os.TempDir(), "non_existent_dir_123456")

	// Create a default config that doesn't exclude anything
	configContent := `
include:
  - "**/*"
exclude: []
`

	configFilePath := filepath.Join(os.TempDir(), ".crev-config.yaml")
	createFile(t, configFilePath, configContent)
	defer os.Remove(configFilePath)

	// Initialize Viper to read the config from TempDir
	viper.Reset()
	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	require.NoError(t, err, "Failed to read config file")

	// Prepare to capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	// No originalDir or directory changing here since path doesn't exist
	rootCmd.SetArgs([]string{"bundle", nonExistentDir})

	err = rootCmd.Execute()
	assert.Error(t, err, "Expected bundle command to fail for non-existent path")

	// Check log messages for appropriate error
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "no files found to bundle", "Should log no files found error")
}
