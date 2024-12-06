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
// Explicit files (provided by --files flag) override any exclude patterns.
func GetAllFilePaths(root string, includePatterns, excludePatterns, explicitFiles []string) ([]string, error) {
	// Normalize root path to absolute path
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	processedExcludePatterns := preprocessExcludePatterns(absRoot, excludePatterns)

	// Handle explicit files: add them to the results and keep track of them
	filePaths, explicitPaths, err := collectExplicitFiles(absRoot, explicitFiles)
	if err != nil {
		return nil, err
	}

	// Now walk the directory and handle non-explicit files
	collectedPaths, err := walkAndCollectPaths(absRoot, includePatterns, processedExcludePatterns, explicitPaths, filePaths)
	if err != nil {
		return nil, err
	}

	return collectedPaths, nil
}

// collectExplicitFiles adds explicit files (those specified by --files) to the output list,
// ensuring they exist and tracking them for later checks.
// This function returns the initial file list (with explicit files) and a map of explicit paths.
func collectExplicitFiles(absRoot string, explicitFiles []string) (filePaths []string, explicitPaths map[string]bool, err error) {
	explicitPaths = make(map[string]bool)

	// First, add explicit files and track their paths
	for _, file := range explicitFiles {
		absPath, err := filepath.Abs(file)
		if err != nil {
			return nil, nil, err
		}
		if _, err := os.Stat(absPath); err == nil {
			explicitPaths[absPath] = true
			filePaths = append(filePaths, absPath)
		}
	}

	return filePaths, explicitPaths, nil
}

// walkAndCollectPaths walks the directory from absRoot, applying exclude patterns, include patterns,
// and considering explicit files. It returns a full list of file paths that meet the criteria.
func walkAndCollectPaths(absRoot string, includePatterns, processedExcludePatterns []string, explicitPaths map[string]bool, initialFiles []string) ([]string, error) {
	filePaths := append([]string(nil), initialFiles...) // copy to avoid mutation
	seenPaths := make(map[string]bool)
	for _, path := range filePaths {
		seenPaths[path] = true
	}

	err := filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == absRoot {
			return nil
		}

		// Skip if we've already seen this path (explicit files)
		if seenPaths[path] {
			return nil
		}

		// Get path relative to root for pattern matching
		relPath, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath) // Convert to forward slashes for consistent pattern matching

		// Determine if this path is excluded and if it's a parent of an explicit file
		excluded, isParentOfExplicit, err := isExcludedPath(absRoot, relPath, processedExcludePatterns, explicitPaths)
		if err != nil {
			return err
		}

		// If this directory (or file) is excluded and not a parent of an explicit file, skip it
		if excluded && !isParentOfExplicit {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// If this is a directory that's excluded but is a parent of an explicit file,
		// we do not add it to the filePaths, but we do continue traversal (do not skip).
		// We only want explicit files, not the directory itself.
		if d.IsDir() && excluded && isParentOfExplicit {
			// Don't add directory to filePaths, just continue walking
			return nil
		}

		// Check include patterns
		include, err := shouldIncludePath(relPath, includePatterns)
		if err != nil {
			return err
		}

		// If we are including this path, add it to the results
		if include {
			filePaths = append(filePaths, path)
			seenPaths[path] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filePaths, nil
}

// isExcludedPath checks if any parent directory of relPath (including itself) matches the exclude patterns.
// It returns whether the path is excluded and whether it is a parent of an explicit file.
//
// If a directory is excluded but also a parent directory of an explicit file, we set isParentOfExplicit = true.
// This allows traversal of the directory without adding it to the output, so that explicit files can be found.
func isExcludedPath(absRoot, relPath string, processedExcludePatterns []string, explicitPaths map[string]bool) (bool, bool, error) {
	dirPath := relPath
	excluded := false
	isParentOfExplicit := false

	for dirPath != "." {
		for _, pattern := range processedExcludePatterns {
			matched, err := doublestar.PathMatch(pattern, dirPath)
			if err != nil {
				return false, false, err
			}
			if matched {
				excluded = true
				// Check if this excluded directory is a parent of any explicit file
				absDir := filepath.Join(absRoot, dirPath)
				for explicit := range explicitPaths {
					if strings.HasPrefix(explicit, absDir+string(os.PathSeparator)) {
						isParentOfExplicit = true
						break
					}
				}
				if isParentOfExplicit {
					// Even though it's excluded, it's a parent of explicit file
					// We'll let traversal continue, but we won't add this directory to filePaths.
					return excluded, isParentOfExplicit, nil
				} else {
					// This directory is excluded and not a parent of any explicit file.
					// We can return now knowing it's excluded without explicit override.
					return excluded, isParentOfExplicit, nil
				}
			}
		}
		dirPath = filepath.Dir(dirPath)
	}

	return excluded, isParentOfExplicit, nil
}

// shouldIncludePath checks whether a path should be included based on the provided includePatterns.
// If no includePatterns are provided, everything is included by default.
func shouldIncludePath(relPath string, includePatterns []string) (bool, error) {
	// Include everything if no patterns specified
	include := len(includePatterns) == 0
	if len(includePatterns) > 0 {
		for _, pattern := range includePatterns {
			matched, err := doublestar.PathMatch(pattern, relPath)
			if err != nil {
				return false, err
			}
			if matched {
				include = true
				break
			}
		}
	}
	return include, nil
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
