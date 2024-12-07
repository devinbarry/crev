package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterEmptyDirectories_NoPaths(t *testing.T) {
	var filePaths []string
	result := filterEmptyDirectories(filePaths)
	require.Empty(t, result, "Expected no output when no input paths are given")
}

func TestFilterEmptyDirectories_AllDirectoriesNoFiles(t *testing.T) {
	// Simulate a structure: rootDir, rootDir/subdir, rootDir/subdir/empty_subdir
	// with no actual files.
	rootDir := t.TempDir()
	subdir := filepath.Join(rootDir, "subdir")
	subSubdir := filepath.Join(subdir, "empty_subdir")

	require.NoError(t, os.MkdirAll(subSubdir, 0755))

	filePaths := []string{rootDir, subdir, subSubdir}
	result := filterEmptyDirectories(filePaths)

	// No directories contain files, so all should be removed, except the root if it's considered a file path.
	// The implementation might consider the root directory as part of the structure if it was passed in.
	// If we only consider directories that were discovered by the walk, we might expect an empty slice.
	// Assuming rootDir counts as a "directory with no files," it should also be removed.
	require.Empty(t, result, "Expected all directories without files to be removed")
}

func TestFilterEmptyDirectories_DirectoriesWithFiles(t *testing.T) {
	// Create a structure:
	// rootDir/
	// ├── file1.go
	// └── subdir/
	//     ├── nested_subdir/
	//     └── file2.txt

	rootDir := t.TempDir()
	file1 := filepath.Join(rootDir, "file1.go")
	subdir := filepath.Join(rootDir, "subdir")
	require.NoError(t, os.MkdirAll(subdir, 0755))
	file2 := filepath.Join(subdir, "file2.txt")

	require.NoError(t, os.WriteFile(file1, []byte("content"), 0644))
	require.NoError(t, os.WriteFile(file2, []byte("content"), 0644))

	// Include the directories and files in filePaths.
	filePaths := []string{
		rootDir,
		file1,
		subdir,
		file2,
	}

	result := filterEmptyDirectories(filePaths)
	// Both rootDir and subdir contain at least one file.
	// No directories should be removed because each has a file (rootDir has file1.go, subdir has file2.txt).
	require.ElementsMatch(t, filePaths, result, "Expected directories with files to remain unchanged")
}

func TestFilterEmptyDirectories_MixedStructure(t *testing.T) {
	// Create a structure:
	// rootDir/
	// ├── file1.go
	// ├── subdir_1/
	// │   ├── file2.go
	// │   └── nested_subdir_1/
	// │       └── file3.go
	// ├── subdir_2/
	// │   └── nested_subdir_2/
	// └── empty_dir/

	rootDir := t.TempDir()

	subdir1 := filepath.Join(rootDir, "subdir_1")
	nested1 := filepath.Join(subdir1, "nested_subdir_1")
	subdir2 := filepath.Join(rootDir, "subdir_2")
	nested2 := filepath.Join(subdir2, "nested_subdir_2")
	emptyDir := filepath.Join(rootDir, "empty_dir")

	require.NoError(t, os.MkdirAll(nested1, 0755))
	require.NoError(t, os.MkdirAll(nested2, 0755))
	require.NoError(t, os.MkdirAll(emptyDir, 0755))

	file1 := filepath.Join(rootDir, "file1.go")
	file2 := filepath.Join(subdir1, "file2.go")
	file3 := filepath.Join(nested1, "file3.go")

	require.NoError(t, os.WriteFile(file1, []byte("content"), 0644))
	require.NoError(t, os.WriteFile(file2, []byte("content"), 0644))
	require.NoError(t, os.WriteFile(file3, []byte("content"), 0644))

	filePaths := []string{
		rootDir,
		file1,
		subdir1,
		file2,
		nested1,
		file3,
		subdir2,
		nested2,
		emptyDir,
	}

	result := filterEmptyDirectories(filePaths)

	// Directories subdir_1 and nested_subdir_1 should remain since they contain files (file2.go, file3.go).
	// rootDir should remain (it has file1.go).
	// subdir_2 and nested_subdir_2 should be removed (no files under them).
	// empty_dir should be removed (no files).
	expected := []string{
		rootDir,
		file1,
		subdir1,
		file2,
		nested1,
		file3,
	}
	require.ElementsMatch(t, expected, result, "Expected only directories containing files or leading to files to remain")
}
