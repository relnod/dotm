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

// ErrorIsNil assert, that the given error is nil.
func ErrorIsNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected error to be nil. Got '%s'", err.Error())
	}
}
