package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var generateCmd = &cobra.Command{
	Use:   "bundle [path]",
	Short: "Bundle your project files into a single file",
	Long: `Bundle your project files into a single file, starting from the specified directory.

File Selection Rules:
1. If --files is specified:
   - Files must exist
   - Listed files are always included, regardless of exclude patterns
   - Additional files can be added via include patterns
   - Exclude patterns still apply to files matched by include patterns

2. If --include patterns are specified (via flags or config):
   - Files matching any include pattern are included
   - Files matching any exclude pattern are excluded (unless specified via --files)
   - Exclude patterns take precedence over include patterns for matched files

3. If neither --files nor --include patterns are specified:
   - Default include pattern "**/*" is used
   - Files matching any exclude pattern are excluded

Config File Integration:
- Values in .crev-config.yaml are used as defaults
- Command line flags override config file values
- Config file include/exclude patterns are merged with command line patterns

Example usage:
  # Use default include pattern (**/*) with default excludes
  crev bundle

  # Bundle specific files (guaranteed inclusion) plus any files matching includes
  crev bundle --files file1.go --files file2.py --include='src/**'

  # Force include a file that would normally be excluded
  crev bundle --files src/vendor/important.go --include='src/**' --exclude='src/vendor/**'

  # Use custom include patterns with default excludes
  crev bundle --include='src/**' --include='lib/**'

  # Combine include and exclude patterns
  crev bundle --include='src/**' --exclude='src/vendor/**'

  # Bundle from a different directory
  crev bundle /path/to/project`,
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

		// Get flags and apply defaults
		explicitFiles := viper.GetStringSlice("files")
		includePatterns := viper.GetStringSlice("include")
		opts.ExcludePatterns = viper.GetStringSlice("exclude")

		// TODO If files are explicitly specified, check that they exist
		if len(explicitFiles) > 0 {
			opts.ExplicitFiles = explicitFiles
		} else {
			// If no files specified, check include patterns
			if len(includePatterns) > 0 {
				opts.IncludePatterns = includePatterns
			} else {
				// If no includes specified, use default include pattern
				opts.IncludePatterns = []string{"**/*"}
			}
		}

		// Execute the bundle operation
		return Bundle(opts)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Add flags without defaults - we'll handle defaults in the RunE function
	generateCmd.Flags().StringSliceP("files", "f", nil,
		"Specify files to always include (overrides exclude patterns for these files)")

	generateCmd.Flags().StringSliceP("include", "i", nil,
		"Include files matching these glob patterns (e.g., 'src/**', '**/*.go')")

	generateCmd.Flags().StringSliceP("exclude", "e", nil,
		"Exclude files matching these glob patterns (except those specified by --files)")

	// Bind flags to viper
	viper.BindPFlag("files", generateCmd.Flags().Lookup("files"))
	viper.BindPFlag("include", generateCmd.Flags().Lookup("include"))
	viper.BindPFlag("exclude", generateCmd.Flags().Lookup("exclude"))
}
