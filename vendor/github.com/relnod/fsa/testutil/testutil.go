package testutil

import (
	"github.com/relnod/fsa"
)

// IsSymlink checks if the given path is a symlink.
func IsSymlink(fs fsa.FileSystem, path string) bool {
	_, err := fs.Stat(path)
	if err != nil {
		return false
	}
	_, err = fs.Readlink(path)
	if err != nil {
		return false
	}
	return true
}

// FileExists checks if the given file exists.
func FileExists(fs fsa.FileSystem, path string) bool {
	_, err := fs.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// DirExists checks if the given directory exists.
func DirExists(fs fsa.FileSystem, path string) bool {
	f, err := fs.Stat(path)
	if err != nil {
		return false
	}
	return f.IsDir()
}
