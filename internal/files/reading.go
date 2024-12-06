package files

import (
	"os"
	"sync"
)

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
