package testutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/relnod/dotm/internal/testutil"
	"github.com/relnod/dotm/internal/testutil/assert"
)

func TestPathIsDir(t *testing.T) {
	var tests = []struct {
		path  testutil.Path
		isDir bool
	}{
		{"a", false},
		{"b/a", false},
		{"b/", true},
	}

	for _, test := range tests {
		if test.path.IsDir() != test.isDir {
			t.Errorf("Expected '%s' to be a directory", test.path)
		}
	}
}

func TestFileSystemPath(t *testing.T) {
	fs := testutil.NewFileSystem()
	basePath := fs.BasePath()

	absolutePath := fs.Path("test")
	expectedPath := filepath.Join(basePath, "test")

	if absolutePath != expectedPath {
		t.Errorf("Expected path to be '%s', but got '%s'", expectedPath, absolutePath)
	}
}

func TestFileSystemJoin(t *testing.T) {
	fs := testutil.NewFileSystem()
	basePath := fs.BasePath()

	absolutePath := fs.Join("test", "test2")
	expectedPath := filepath.Join(basePath, "test/test2")

	if absolutePath != expectedPath {
		t.Errorf("Expected path to be '%s', but got '%s'", expectedPath, absolutePath)
	}
}

func TestFileSystemCleanup(t *testing.T) {
	fs := testutil.NewFileSystem()
	basePath := fs.BasePath()

	os.MkdirAll(filepath.Join(basePath, "test"), os.ModePerm)

	err := fs.Cleanup()
	assert.ErrorIsNil(t, err)

	if _, err := os.Stat(fs.BasePath()); err == nil {
		t.Errorf("Expected basePath to not exist")
	}
}

func TestFileSystemMkdirAll(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{"test"},
		{"a/b"},
	}

	for _, test := range tests {
		t.Run(test.path, func(tt *testing.T) {
			fs := testutil.NewFileSystem()

			err := fs.MkdirAll(test.path)
			defer fs.Cleanup()

			assert.ErrorIsNil(tt, err)
			assert.PathExists(tt, fs.Path(test.path), true)
		})
	}
}

func TestFileSystemCreate(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{"test"},
		{"foo/bar"},
		{"a/b"},
	}

	for _, test := range tests {
		t.Run(test.path, func(tt *testing.T) {
			fs := testutil.NewFileSystem()

			err := fs.Create(test.path)
			defer fs.Cleanup()

			assert.ErrorIsNil(tt, err)
			assert.PathExists(tt, fs.Path(test.path), false)
		})
	}
}

func TestFileSystemFromFileStructure(t *testing.T) {
	fs := testutil.NewFileSystem()

	fileStructure := testutil.FileStructure{
		"test/",
		"test/a",
		"a/b",
		"c",
	}

	var expected = []struct {
		path  string
		isDir bool
	}{
		{"test", true},
		{"test/a", false},
		{"a/b", false},
		{"c", false},
	}

	err := fs.CreateFromFileStructure(fileStructure)
	defer fs.Cleanup()

	assert.ErrorIsNil(t, err)

	for _, exp := range expected {
		assert.PathExists(t, fs.Path(exp.path), exp.isDir)
	}

}
