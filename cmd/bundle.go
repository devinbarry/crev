// Description: This file contains the generate command which generates a textual representation of the project structure.
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/devinbarry/crev/internal/files"
	"github.com/devinbarry/crev/internal/formatting"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Specific prefixes to ignore
var specificPrefixesToIgnore = []string{
	"crev", // ignore crev specific files
}

// Specific extensions to ignore
var specificExtensionsToIgnore = []string{
	".jpeg",
	".jpg",
	".png",
	".gif",
	".pdf",
	".svg",
	".ico",
	".woff",
	".woff2",
	".eot",
	".ttf",
	".otf",
}

// Specific filenames to ignore
var specificFilesToIgnore = []string{
	"Thumbs.db",
	"poetry.lock",
	"__pycache__",
}

var generateCmd = &cobra.Command{
	Use:   "bundle [path]",
	Short: "Bundle your project into a single file",
	Long: `Bundle your project into a single file, starting from the specified directory.
By default, all files are included unless they match an exclude pattern.

Use the --include and --exclude flags to specify patterns for files and directories to include or exclude.
Patterns are processed in order, and files matching any exclude pattern will be excluded even if they match an include pattern.
Use the -f or --files flag to specify explicit files to include, overriding exclude patterns if necessary.

Example usage:
  crev bundle
  crev bundle /path/to/project
  crev bundle --exclude='*.md' --exclude='test/*'
  crev bundle --include='src/**' --exclude='src/vendor/**'
  crev bundle -f file1.go,file2.py,file3.md
`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Start timer
		start := time.Now()

		// Get current working directory for output file path
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		rootDir := "."
		if len(args) > 0 {
			rootDir = args[0]
		}

		// Check if rootDir exists before proceeding
		if _, err := os.Stat(rootDir); os.IsNotExist(err) {
			errMsg := "no files found to bundle"
			log.Print(errMsg)
			return fmt.Errorf("%s. Please check your include/exclude patterns and the specified path", errMsg)
		}

		// Get flags
		explicitFiles := viper.GetStringSlice("files")
		includePatterns := viper.GetStringSlice("include")
		excludePatterns := viper.GetStringSlice("exclude")

		// Add excludes for prefixes
		for _, prefix := range specificPrefixesToIgnore {
			// Exclude directories and files starting with the prefix at any level
			excludePatterns = append(excludePatterns, "**/"+prefix+"*", prefix+"*")
		}

		// Convert extensions to exclude patterns
		for _, ext := range specificExtensionsToIgnore {
			excludePatterns = append(excludePatterns, "**/*"+ext)
		}

		// Add specific filenames to exclude patterns
		for _, file := range specificFilesToIgnore {
			excludePatterns = append(excludePatterns, "**/"+file)
		}

		// Create output file in current working directory
		outputFile := filepath.Join(cwd, "crev-project.txt")

		// Fetch file paths
		filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, explicitFiles)
		if err != nil {
			return fmt.Errorf("error getting file paths: %w", err)
		}

		if len(filePaths) == 0 {
			errMsg := "no files found to bundle"
			log.Print(errMsg)
			return fmt.Errorf("%s. Please check your include/exclude patterns and the specified path", errMsg)
		}

		// Generate the project tree (structure)
		projectTree := formatting.GeneratePathTree(filePaths)
		maxConcurrency := 100

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

		// Log success
		log.Printf("Project overview successfully saved to: %s", outputFile)
		log.Printf("Estimated token count: %d - %d tokens", len(projectString)/4, len(projectString)/3)
		log.Printf("Execution time: %s", time.Since(start))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Add the -f flag (short for --files) to specify explicit files to include in the bundle
	generateCmd.Flags().StringSliceP("files", "f", []string{}, "Specify multiple file paths to include (e.g., --files file1.go --files file2.py)")

	// Add the --include flag to include files or directories matching patterns
	generateCmd.Flags().StringSliceP("include", "I", []string{}, "Include files or directories matching these glob patterns (e.g., 'src/**', '**/*.go')")

	// Add the --exclude flag to exclude files or directories matching patterns
	generateCmd.Flags().StringSliceP("exclude", "E", []string{}, "Exclude files or directories matching these glob patterns (e.g., 'vendor/**', '**/*.test.go')")

	// Bind flags to viper for easy retrieval
	viper.BindPFlag("files", generateCmd.Flags().Lookup("files"))
	viper.BindPFlag("include", generateCmd.Flags().Lookup("include"))
	viper.BindPFlag("exclude", generateCmd.Flags().Lookup("exclude"))
}
