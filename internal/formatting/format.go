// Contains code to format the project structure into a string.
package formatting

import (
	"path/filepath"
	"sort"
	"strings"
)

// node represents a node in a tree structure. Each node has a name (which could
// be a directory or file name) and a map of children, where the key is the child
// name and the value is a pointer to another node. This structure is used to
// build and represent hierarchical directory paths, with each node corresponding
// to a folder or file in the directory tree.
type node struct {
	name     string
	children map[string]*node
}

// GeneratePathTree Given a list of paths, GeneratePathTree returns a string representation of the
// directory structure.
func GeneratePathTree(paths []string) string {
	root := &node{children: make(map[string]*node)}

	// Sort the paths lexicographically to ensure correct tree structure
	sort.Strings(paths)

	// Build the tree structure
	for _, path := range paths {
		cleanedPath := filepath.Clean(path)
		if cleanedPath == "." {
			continue // Skip if the path is empty or root
		}
		parts := strings.Split(filepath.ToSlash(cleanedPath), "/")
		current := root
		for _, part := range parts {
			if _, exists := current.children[part]; !exists {
				current.children[part] = &node{name: part, children: make(map[string]*node)}
			}
			current = current.children[part]
		}
	}

	// Generate the tree string
	var sb strings.Builder
	printTree(root, "", &sb)
	return sb.String()
}

func printTree(n *node, prefix string, sb *strings.Builder) {
	children := make([]*node, 0, len(n.children))
	for _, child := range n.children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].name < children[j].name
	})

	for i, child := range children {
		isLast := i == len(children)-1
		sb.WriteString(prefix)
		if isLast {
			sb.WriteString("└── ")
		} else {
			sb.WriteString("├── ")
		}
		sb.WriteString(child.name)
		sb.WriteString("\n") // Always append a newline

		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		printTree(child, newPrefix, sb)
	}
}

// CreateProjectString Creates a string representation of the project.
func CreateProjectString(projectTree string, fileContentMap map[string]string) string {
	var projectString strings.Builder
	projectString.WriteString("Project Directory Structure:" + "\n")
	projectString.WriteString(projectTree + "\n\n")

	// Collect and sort the file paths lexicographically to make the function deterministic
	filePaths := make([]string, 0, len(fileContentMap))
	for filePath := range fileContentMap {
		filePaths = append(filePaths, filePath)
	}
	sort.Strings(filePaths)

	for _, fileName := range filePaths {
		fileContent := fileContentMap[fileName]
		// Skip displaying the file if it has no content
		if strings.TrimSpace(fileContent) == "" {
			continue
		}
		// Add file name and content if the file has non-empty content
		projectString.WriteString("File: " + "\n")
		projectString.WriteString(fileName + "\n")
		projectString.WriteString("Content: " + "\n")
		projectString.WriteString(fileContent + "\n\n")
	}
	return projectString.String()
}
