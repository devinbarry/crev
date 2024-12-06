package files

import (
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// GetAllFilePaths returns all the file paths in the root directory and its subdirectories,
// while respecting inclusion and exclusion patterns.
// After collecting files from walking the directory and applying include/exclude patterns,
// explicit files provided by the user with the --files flag are added. This ensures that
// explicitly specified files (via --files) override any exclude patterns.
func GetAllFilePaths(root string, includePatterns, excludePatterns, explicitFiles []string) ([]string, error) {
	// Normalize root path to absolute path
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	processedExcludePatterns := preprocessExcludePatterns(absRoot, excludePatterns)

	// First, handle explicit files as they should override exclude patterns
	explicitFilePaths := make(map[string]bool)
	for _, file := range explicitFiles {
		absPath, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(absPath); err == nil {
			explicitFilePaths[absPath] = true
			filePaths = append(filePaths, absPath)
		}
	}

	// Walk the directory
	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == absRoot {
			return nil
		}

		// If this is an explicit file, we've already handled it
		if explicitFilePaths[path] {
			return nil
		}

		// Get path relative to root for pattern matching
		relPath, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}

		// Convert to forward slashes for consistent pattern matching
		relPath = filepath.ToSlash(relPath)

		// Check exclude patterns first
		for _, pattern := range processedExcludePatterns {
			matched, err := doublestar.PathMatch(pattern, relPath)
			if err != nil {
				return err
			}
			if matched {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Check include patterns
		include := len(includePatterns) == 0 // Include everything if no patterns specified
		if len(includePatterns) > 0 {
			for _, pattern := range includePatterns {
				matched, err := doublestar.PathMatch(pattern, relPath)
				if err != nil {
					return err
				}
				if matched {
					include = true
					break
				}
			}
		}

		if include {
			filePaths = append(filePaths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filePaths, nil
}

// preprocessExcludePatterns adjusts exclude patterns to handle directories and trailing slashes.
// For directories, it adds both the directory itself and "/**" pattern to exclude all contents.
// For files or non-existent paths, it uses the pattern as-is.
// Empty patterns are skipped to avoid unintended matches.
func preprocessExcludePatterns(root string, excludePatterns []string) []string {
	var processedPatterns []string

	for _, pattern := range excludePatterns {
		// Skip empty patterns
		if pattern == "" {
			continue
		}

		// Clean the pattern by removing trailing slashes
		cleanPattern := strings.TrimRight(pattern, "/\\")

		// Check if the pattern corresponds to an existing path
		fullPath := filepath.Join(root, cleanPattern)
		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			// For directories, add both the directory pattern and its contents
			processedPatterns = append(processedPatterns,
				cleanPattern,       // Match the directory itself
				cleanPattern+"/**", // Match all contents
			)
		} else {
			// For files or non-existent paths, use the cleaned pattern
			processedPatterns = append(processedPatterns, cleanPattern)
		}
	}

	return processedPatterns
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
