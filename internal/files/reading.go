package files

import (
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// GetAllFilePaths returns all the file paths in the root directory and its subdirectories,
// while respecting inclusion and exclusion patterns.
// After collecting files from walking the directory and applying include/exclude patterns,
// explicit files provided by the user with the --files flag are added. This ensures that
// explicitly specified files (via --files) override any exclude patterns.
//
// This function returns all paths as absolute paths to maintain consistency with tests that
// expect absolute paths.
func GetAllFilePaths(root string, includePatterns, excludePatterns, explicitFiles []string) ([]string, error) {
	// Canonicalize the root directory to avoid symlink issues (e.g., /var vs /private/var)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	absRoot, err = filepath.EvalSymlinks(absRoot)
	if err != nil {
		return nil, err
	}

	var filePaths []string

	processedExcludePatterns := preprocessExcludePatterns(absRoot, excludePatterns)

	// Walk the directory using the canonical root
	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute the relative path from the canonical root for pattern matching
		relPath, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Check exclude patterns
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
		include := len(includePatterns) == 0
		if len(includePatterns) > 0 {
			include = false
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
			// Canonicalize the path to ensure consistency
			canonicalPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			filePaths = append(filePaths, canonicalPath)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Add explicit files after processing include/exclude
	// Explicit files override exclude patterns.
	for _, file := range explicitFiles {
		absPath, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}
		absPath, err = filepath.EvalSymlinks(absPath)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(absPath); err == nil {
			if !contains(filePaths, absPath) {
				filePaths = append(filePaths, absPath)
			}
		}
	}

	return filePaths, nil
}

// preprocessExcludePatterns adjusts exclude patterns to handle directories and trailing slashes
func preprocessExcludePatterns(root string, excludePatterns []string) []string {
	var processedPatterns []string

	for _, pattern := range excludePatterns {
		adjustedPattern := pattern

		// Remove trailing slashes for consistency
		adjustedPattern = strings.TrimSuffix(adjustedPattern, string(os.PathSeparator))

		// Check if the pattern corresponds to a directory
		dirPath := filepath.Join(root, adjustedPattern)
		if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
			// Append /** to match all contents within the directory
			adjustedPattern = filepath.ToSlash(filepath.Clean(adjustedPattern)) + "/**"
		}

		// Add both the directory and its contents to the patterns
		processedPatterns = append(processedPatterns, adjustedPattern)
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

// getFileContent returns the content of the given file.
func getFileContent(filePath string) (string, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// GetContentMapOfFiles returns a map of file paths to their content.
func GetContentMapOfFiles(filePaths []string, maxConcurrency int) (map[string]string, error) {
	var fileContentMap sync.Map
	var wg sync.WaitGroup
	errChan := make(chan error, len(filePaths))
	semaphore := make(chan struct{}, maxConcurrency)

	for _, path := range filePaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			info, err := os.Stat(p)
			if err != nil {
				errChan <- err
				return
			}
			if !info.IsDir() {
				fileContent, err := getFileContent(p)
				if err != nil {
					errChan <- err
					return
				}
				fileContentMap.Store(p, fileContent)
			} else {
				dirEntries, err := os.ReadDir(p)
				if err != nil {
					errChan <- err
					return
				}
				if len(dirEntries) == 0 {
					fileContentMap.Store(p, "empty directory")
				}
			}
		}(path)
	}
	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	resultMap := make(map[string]string)
	fileContentMap.Range(func(key, value interface{}) bool {
		resultMap[key.(string)] = value.(string)
		return true
	})

	return resultMap, nil
}
