package cmd

import (
	"fmt"
	"github.com/devinbarry/crev/internal/files"
	"github.com/devinbarry/crev/internal/formatting"
	"log"
	"os"
	"path/filepath"
	"time"
)

// BundleOptions contains all the configuration options for the bundle operation
type BundleOptions struct {
	RootDir         string
	ExplicitFiles   []string
	IncludePatterns []string
	ExcludePatterns []string
	OutputDir       string
	MaxConcurrency  int
}

// DefaultBundleOptions returns a BundleOptions with default values
func DefaultBundleOptions() BundleOptions {
	return BundleOptions{
		RootDir:        ".",
		MaxConcurrency: 100,
	}
}

// validateExplicitFiles checks if all explicitly specified files exist
func validateExplicitFiles(files []string) error {
	var missingFiles []string
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(missingFiles) > 0 {
		return fmt.Errorf("the following files specified via --files do not exist: %v", missingFiles)
	}
	return nil
}

// Bundle performs the main bundling operation
func Bundle(opts BundleOptions) error {
	start := time.Now()
	log.Printf("Starting bundle operation in directory: %s", opts.RootDir)

	// Get absolute path for better error messaging
	absRootDir, err := filepath.Abs(opts.RootDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path %q: %w", opts.RootDir, err)
	}

	// Check if rootDir exists and is accessible
	if _, err := os.Stat(absRootDir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %q does not exist", absRootDir)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied accessing directory %q", absRootDir)
		}
		return fmt.Errorf("error accessing directory %q: %w", absRootDir, err)
	}

	// Validate explicit files if any are specified
	if len(opts.ExplicitFiles) > 0 {
		if err := validateExplicitFiles(opts.ExplicitFiles); err != nil {
			return err
		}
	}

	// Add default exclude patterns
	opts.ExcludePatterns = appendDefaultExcludes(opts.ExcludePatterns)
	log.Printf("Files: %v", opts.ExplicitFiles)
	log.Printf("Includes: %v", opts.IncludePatterns)
	log.Printf("Excludes: %v", opts.ExcludePatterns)

	// Create output file path
	outputFile := filepath.Join(opts.OutputDir, "crev-project.txt")

	// Fetch file paths
	filePaths, err := files.GetAllFilePaths(opts.RootDir, opts.IncludePatterns, opts.ExcludePatterns, opts.ExplicitFiles)
	if err != nil {
		return fmt.Errorf("error getting file paths: %w", err)
	}

	log.Println(filePaths)

	if len(filePaths) == 0 {
		return fmt.Errorf("no files found to bundle. Please check your include/exclude patterns and the specified path")
	}

	// Generate and save the bundle
	if err := generateBundle(filePaths, outputFile, opts.MaxConcurrency); err != nil {
		return err
	}

	// Log success
	log.Printf("Project overview successfully saved to: %s", outputFile)
	log.Printf("Execution time: %s", time.Since(start))

	return nil
}

// appendDefaultExcludes adds the default exclude patterns to the provided patterns
func appendDefaultExcludes(patterns []string) []string {
	// Add excludes for prefixes
	for _, prefix := range specificPrefixesToIgnore {
		patterns = append(patterns, "**/"+prefix+"*", prefix+"*")
	}

	// Convert extensions to exclude patterns
	for _, ext := range specificExtensionsToIgnore {
		patterns = append(patterns, "**/*"+ext)
	}

	// Add specific filenames to exclude patterns
	for _, file := range specificFilesToIgnore {
		patterns = append(patterns, "**/"+file)
	}

	return patterns
}

// generateBundle creates the bundle file from the given file paths
func generateBundle(filePaths []string, outputFile string, maxConcurrency int) error {
	// Generate the project tree (structure)
	projectTree := formatting.GeneratePathTree(filePaths)

	// Retrieve file contents
	fileContentMap, err := files.GetContentMapOfFiles(filePaths, maxConcurrency)
	if err != nil {
		return fmt.Errorf("error getting file contents: %w", err)
	}

	// Create and save the project string
	projectString := formatting.CreateProjectString(projectTree, fileContentMap)
	if err := files.SaveStringToFile(projectString, outputFile); err != nil {
		return fmt.Errorf("error saving file: %w", err)
	}

	log.Printf("Estimated token count: %d - %d tokens", len(projectString)/4, len(projectString)/3)
	return nil
}
