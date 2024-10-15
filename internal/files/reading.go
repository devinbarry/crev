// Contains code to read the content of files and directories.
package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// GetExplicitFilePaths validates and returns the file paths from the given list that exist.
func GetExplicitFilePaths(explicitPaths []string) ([]string, error) {
	var validPaths []string

	// Iterate through the provided paths and check if they exist.
	for _, path := range explicitPaths {
		// Check if the file or directory exists
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			return nil, err // Return an error if the file or directory does not exist
		} else if err != nil {
			return nil, err // Handle other potential errors with os.Stat
		}

		// If it exists, add it to the validPaths list
		// Only add file paths (excluding directories), but directories can be added if required.
		if !info.IsDir() {
			validPaths = append(validPaths, path)
		}
	}

	return validPaths, nil
}

// GetAllFilePaths Given a root path returns all the file paths in the root directory
// and its subdirectories, while respecting exclusion rules.
func GetAllFilePaths(root string, prefixesToFilter []string, extensionsToKeep []string,
	extensionsToIgnore []string, excludeList []string) ([]string, error) {

	var filePaths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory.
		if path == root {
			return nil
		}
		// First filter out the paths that contain any of the prefixes in prefixesToFilter.
		for _, prefixToFilter := range prefixesToFilter {
			if strings.HasPrefix(filepath.Base(path), prefixToFilter) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		// Filter out paths in the exclude list.
		for _, excludePath := range excludeList {
			if strings.HasPrefix(path, excludePath) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		// Filter out the files that have the extensions in extensionsToIgnore.
		for _, ext := range extensionsToIgnore {
			if filepath.Ext(path) == ext {
				return nil
			}
		}
		// Process file based on extension filters.
		if d.IsDir() || len(extensionsToKeep) == 0 {
			filePaths = append(filePaths, path)
			return nil
		}
		for _, ext := range extensionsToKeep {
			if filepath.Ext(path) == ext {
				filePaths = append(filePaths, path)
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return filePaths, nil
}

// Given a file path, GetFileContent returns the content of the file.
func getFileContent(filePath string) (string, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// GetContentMapOfFiles Given a list of file paths, returns a map of file paths to their content.
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
