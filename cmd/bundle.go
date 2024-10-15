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
	"package",
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
	Use:   "bundle",
	Short: "Bundle your project into a single file",
	Long: `Bundle your project into a single file, starting from the directory you are in.
By default, only files explicitly specified via --include-ext will have their contents included in the bundle, while the directory structure will be preserved for all files.
Use the --all flag to include all file contents (excluding ignored ones).
For more information see: https://crevcli.com/docs

Example usage:
crev bundle
crev bundle --ignore-pre=tests,readme --ignore-ext=.txt 
crev bundle --ignore-pre=tests,readme --include-ext=.go,.py,.js
crev bundle --all
`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		// start timer
		start := time.Now()

		// get all file paths from the root directory
		rootDir := "."

		// Get the --all flag
		bundleAll := viper.GetBool("all")

		// Get prefixes and extensions to ignore/include
		prefixesToIgnore := viper.GetStringSlice("ignore-pre")
		prefixesToIgnore = append(prefixesToIgnore, standardPrefixesToIgnore...)

		extensionsToIgnore := viper.GetStringSlice("ignore-ext")
		extensionsToIgnore = append(extensionsToIgnore, standardExtensionsToIgnore...)

		extensionsToInclude := viper.GetStringSlice("include-ext")

		// Fetch file paths
		filePaths, err := files.GetAllFilePaths(rootDir, prefixesToIgnore,
			extensionsToInclude, extensionsToIgnore)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Generate the project tree (structure) regardless of the --all flag
		projectTree := formatting.GeneratePathTree(filePaths)

		maxConcurrency := 100
		// Fetch the file contents based on the --all flag
		var fileContentMap map[string]string
		if bundleAll {
			// Include all file contents
			fileContentMap, err = files.GetContentMapOfFiles(filePaths, maxConcurrency)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Include only the file contents for files that match --include-ext
			if len(extensionsToInclude) > 0 {
				fileContentMap, err = files.GetContentMapOfFiles(filePaths, maxConcurrency)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				// If no extensions to include were specified, create an empty map
				fileContentMap = make(map[string]string)
			}
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

	// Add the --all flag to include all files' content
	generateCmd.Flags().Bool("all", false, "Include all file contents in the bundle (excluding ignored files)")

	// Add existing flags for ignoring and including extensions/prefixes
	generateCmd.Flags().StringSlice("ignore-pre", []string{}, "Comma-separated prefixes of file and dir names to ignore. Ex tests,readme")
	generateCmd.Flags().StringSlice("ignore-ext", []string{}, "Comma-separated file extensions to ignore. Ex .txt,.md")
	generateCmd.Flags().StringSlice("include-ext", []string{}, "Comma-separated file extensions to include. Ex .go,.py,.js")

	// Bind flags to viper for easy retrieval
	err := viper.BindPFlag("all", generateCmd.Flags().Lookup("all"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("ignore-pre", generateCmd.Flags().Lookup("ignore-pre"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("ignore-ext", generateCmd.Flags().Lookup("ignore-ext"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("include-ext", generateCmd.Flags().Lookup("include-ext"))
	if err != nil {
		log.Fatal(err)
	}
}
