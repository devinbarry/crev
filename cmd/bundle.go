// Description: This file contains the generate command which generates a textual representation of the project structure.
package cmd

import (
	"log"
	"time"

	"github.com/devinbarry/crev/internal/files"
	"github.com/devinbarry/crev/internal/formatting"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Standard prefixes to ignore
var standardPrefixesToIgnore = []string{
	// ignore .git, .idea, .vscode, etc.
	".",

	// General project files
	"license",
	"LICENSE",
	"DEVELOP.md",
	"readme",
	"README",
	// ignore crev specific files
	"crev",
	// ignore go.mod, go.sum, etc.
	"go",
	// poetry
	"pyproject.toml",
	"poetry.lock",
	"venv",
	// output files
	"build",
	"dist",
	"out",
	"target",
	"bin",
	// javascript
	"node_modules",
	"coverage",
	"public",
	"static",
	"Thumbs.db",
	"package",
	"yarn.lock",
	"tsconfig",
	// next.js
	"next.config",
	"next-env",

	// python
	"requirements.txt",
	"__pycache__",
	"logs",
	// java
	"gradle",
	// c++
	"CMakeLists",
	// ruby
	"vendor",
	"Gemfile",
	// php
	"composer",
	// tailwind
	"tailwind",
	"postcss",
}

// Standard extensions to ignore
var standardExtensionsToIgnore = []string{
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
	Run: func(_ *cobra.Command, args []string) {
		// Start timer
		start := time.Now()

		rootDir := "."
		if len(args) > 0 {
			rootDir = args[0]
		}

		// Get flags
		explicitFiles := viper.GetStringSlice("files")
		includePatterns := viper.GetStringSlice("include")
		excludePatterns := viper.GetStringSlice("exclude")

		// Incorporate standard prefixes and extensions into exclude patterns
		// Convert prefixes to exclude patterns
		for _, prefix := range standardPrefixesToIgnore {
			// Exclude directories and files starting with the prefix at any level
			excludePatterns = append(excludePatterns, "**/"+prefix+"*", prefix+"*")
		}

		// Convert extensions to exclude patterns
		for _, ext := range standardExtensionsToIgnore {
			excludePatterns = append(excludePatterns, "**/*"+ext)
		}

		// Fetch file paths
		filePaths, err := files.GetAllFilePaths(rootDir, includePatterns, excludePatterns, explicitFiles)
		if err != nil {
			log.Fatal(err)
			return
		}

		if len(filePaths) == 0 {
			log.Fatal("No files found to bundle. Please check your include/exclude patterns and the specified path.")
		}

		// Generate the project tree (structure)
		projectTree := formatting.GeneratePathTree(filePaths)
		maxConcurrency := 100

		// Retrieve file contents
		fileContentMap, err := files.GetContentMapOfFiles(filePaths, maxConcurrency)
		if err != nil {
			log.Fatal(err)
		}

		// Create the project string
		projectString := formatting.CreateProjectString(projectTree, fileContentMap)

		outputFile := "crev-project.txt"
		// Save the project string to a file
		err = files.SaveStringToFile(projectString, outputFile)
		if err != nil {
			log.Fatal(err)
		}

		// Log success
		log.Println("Project overview successfully saved to: " + outputFile)

		// Estimate number of tokens
		log.Printf("Estimated token count: %d - %d tokens",
			len(projectString)/4, len(projectString)/3)

		elapsed := time.Since(start)
		log.Printf("Execution time: %s", elapsed)
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
