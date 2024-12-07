package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// Specific prefixes to ignore
var specificPrefixesToIgnore = []string{
	// ignore .git, .idea, .vscode, etc.
	".",
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
	"go.mod",
	"go.sum",
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
		// Get current working directory for output file path
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// Create bundle options
		opts := DefaultBundleOptions()

		// Set root directory
		if len(args) > 0 {
			opts.RootDir = args[0]
		}

		// Set output directory
		opts.OutputDir = cwd

		// Get flags from viper
		opts.ExplicitFiles = viper.GetStringSlice("files")
		opts.IncludePatterns = viper.GetStringSlice("include")
		opts.ExcludePatterns = viper.GetStringSlice("exclude")

		// Execute the bundle operation
		return Bundle(opts)
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
