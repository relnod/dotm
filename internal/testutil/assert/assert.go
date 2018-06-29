package assert

import (
	"os"
	"testing"
)

// PathExists asserts, that the given path exists and is either a file or a
// directory, depending on isDir.
func PathExists(t *testing.T, path string, isDir bool) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Expected '%s' to exist", path)
	}
	if fileInfo.IsDir() == true && isDir == false {
		t.Fatalf("Expected '%s' to be a directory", path)
	}
	if fileInfo.IsDir() == false && isDir == true {
		t.Fatalf("Expected '%s' to be a file", path)
	}
}

// PathNotExists checks if the path does not exist.
func PathNotExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	if err == nil {
		t.Fatalf("Expected '%s' to NOT exist", path)
	}
}

// IsSymlink checks if the given path is a symlink.
func IsSymlink(t *testing.T, path string) {
	f, err := os.Stat(path)
	if err != nil {
		t.Errorf("Expected '%s' to be a symlink. But the path does not exists!", path)
		return
	}

	_, err = os.Readlink(path)
	if err != nil {
		t.Errorf("Expected '%s' to be a symlink. Got the following modes: '%s'.", path, f.Mode())
	}
}

// ErrorIsNil asserts, that the given error is nil.
func ErrorIsNil(t *testing.T, err error) bool {
	if err != nil {
		t.Errorf("Expected error to be nil. Got '%s'", err.Error())
		return false
	}

	return true
}

// ErrorEquals asserts, that both errors are the same.
func ErrorEquals(t *testing.T, err error, expectedErr error) {
	if expectedErr == nil {
		ErrorIsNil(t, err)
	} else if err != expectedErr {
		t.Errorf("Expected error to equal '%s'. Got '%s'", expectedErr.Error(), err.Error())
	}
}
