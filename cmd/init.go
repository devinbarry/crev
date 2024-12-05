// Description: This file implements the "init" command, which generates a default configuration file in the current directory.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Define a default template configuration
var defaultConfig = []byte(`# Configuration for the crev tool

# Specify the glob patterns for files and directories to include (default is all files)
include:
  - "**/*"

# Specify the glob patterns for files and directories to exclude
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

# Example:
# include:
#   - "src/**"
#   - "**/*.go"
# exclude:
#   - "vendor/**"
#   - "**/*.test.go"
`)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a default configuration file",
	Long: `Generates a default configuration file (.crev-config.yaml) in the current directory.

The configuration file includes:
- Include and exclude patterns for files and directories when generating the project overview.

You can modify this file as needed to suit your project's structure.
`,
	Run: func(cmd *cobra.Command, args []string) {
		configFileName := ".crev-config.yaml"

		// Check if the config file already exists
		if _, err := os.Stat(configFileName); err == nil {
			fmt.Println("Config file already exists at", configFileName)
			os.Exit(1)
		}

		// Write the default config
		err := os.WriteFile(configFileName, defaultConfig, 0644)
		if err != nil {
			fmt.Println("Unable to write config file:", err)
			os.Exit(1)
		}

		// Inform the user
		fmt.Println("Config file created at:", configFileName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
