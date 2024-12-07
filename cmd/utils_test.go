package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// multiWriter writes to multiple io.Writers
type multiWriter struct {
	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
	}
	return len(p), nil
}

// testEnv represents a test environment with temporary directory and logging
type testEnv struct {
	t           *testing.T
	TempDir     string
	OriginalDir string
	LogBuffer   *bytes.Buffer
}

// newTestEnv creates a new test environment with temporary directory and logging setup
func newTestEnv(t *testing.T) *testEnv {
	// Reset viper completely
	viper.Reset()

	// Reset the command's flags
	generateCmd.ResetFlags()

	// Re-add the original flags
	generateCmd.Flags().StringSliceP("files", "f", nil,
		"Specify files to always include (overrides exclude patterns for these files)")
	generateCmd.Flags().StringSliceP("include", "i", nil,
		"Include files matching these glob patterns (e.g., 'src/**', '**/*.go')")
	generateCmd.Flags().StringSliceP("exclude", "e", nil,
		"Exclude files matching these glob patterns (except those specified by --files)")

	// Re-bind flags to viper
	viper.BindPFlag("files", generateCmd.Flags().Lookup("files"))
	viper.BindPFlag("include", generateCmd.Flags().Lookup("include"))
	viper.BindPFlag("exclude", generateCmd.Flags().Lookup("exclude"))

	// Create temporary directory
	tempDir := t.TempDir()

	// Get original directory
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")

	// Setup log buffer
	logBuf := &bytes.Buffer{}

	// Create a multi-writer that writes to both buffer and stderr if verbose testing is enabled
	writers := []io.Writer{logBuf}
	if testing.Verbose() {
		writers = append(writers, os.Stderr)
	}
	mw := &multiWriter{writers: writers}
	log.SetOutput(mw)

	// Change to temp directory
	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change to temporary directory")

	// Setup cleanup
	t.Cleanup(func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err, "Failed to change back to original directory")
		log.SetOutput(os.Stderr)
	})

	return &testEnv{
		t:           t,
		TempDir:     tempDir,
		OriginalDir: originalDir,
		LogBuffer:   logBuf,
	}
}

// setupConfig creates a .crev-config.yaml file with given content and initializes viper
func (env *testEnv) setupConfig(configContent string) {
	configPath := filepath.Join(env.TempDir, ".crev-config.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(env.t, err, "Failed to create config file")

	viper.Reset()
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	require.NoError(env.t, err, "Failed to read config file")
}

// createProjectStructure creates a directory structure from a map of file paths to contents
func (env *testEnv) createProjectStructure(files map[string]string) {
	for path, content := range files {
		fullPath := filepath.Join(env.TempDir, path)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		require.NoError(env.t, err, "Failed to create directory: %s", dir)

		// Create file with content
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(env.t, err, "Failed to create file: %s", path)
	}
}

// assertFileContents checks if the output file contains or doesn't contain expected content
func (env *testEnv) assertFileContents(outputFile string, expectedContent, unexpectedContent []string) {
	_, err := os.Stat(outputFile)
	require.NoError(env.t, err, "Expected output file %s to exist", outputFile)

	content, err := os.ReadFile(outputFile)
	require.NoError(env.t, err, "Failed to read output file")

	contentStr := string(content)
	for _, expected := range expectedContent {
		require.Contains(env.t, contentStr, expected, "%s should be included", expected)
	}
	for _, unexpected := range unexpectedContent {
		require.NotContains(env.t, contentStr, unexpected, "%s should not be included", unexpected)
	}
}

// executeBundleCmd executes the bundle command with given arguments
func (env *testEnv) executeBundleCmd(args ...string) error {
	fullArgs := append([]string{"bundle"}, args...)
	log.Printf("Test executing bundle command with args: %#v", fullArgs)
	rootCmd.SetArgs(fullArgs)
	return rootCmd.Execute()
}

// assertLogContains checks if the log buffer contains expected messages
func (env *testEnv) assertLogContains(expectedMessages ...string) {
	logOutput := env.LogBuffer.String()
	for _, msg := range expectedMessages {
		require.Contains(env.t, logOutput, msg, "Log should contain: %s", msg)
	}
}

// assertErrorContains checks if the error contains expected message
func (env *testEnv) assertErrorContains(err error, expectedMsg string) {
	require.Error(env.t, err, "Expected an error")
	require.Contains(env.t, err.Error(), expectedMsg,
		"Error should contain message: %s", expectedMsg)
}

// assertNoErrorWithLog checks there's no error and log contains message
func (env *testEnv) assertNoErrorWithLog(err error, expectedLog string) {
	require.NoError(env.t, err, "Expected no error")
	require.Contains(env.t, env.LogBuffer.String(), expectedLog,
		"Log should contain: %s", expectedLog)
}

// Common config templates
const (
	basicConfig = `
include:
  - "**/*"
exclude: []
`

	excludeConfig = `
include:
  - "**/*"
exclude:
  - "*.md"
  - "node_modules/**"
  - ".git/**"
  - "*.png"
`

	fullConfig = `
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
)
