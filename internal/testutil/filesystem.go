package testutil

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// FileStructure is a representation of a file structure. Each entry represent
// the full path of a file or directory.
type FileStructure []Path

// Path represents a file path.
type Path string

// IsDir returns true if the last character of the file path is a "/".
func (p Path) IsDir() bool {
	if p[len(p)-1:] == "/" {
		return true
	}

	return false
}

// FileSystem provides an easy way to create a temporary file and directory
// structure.
// TODO: maybe there is a better word, than "FileSystem"
type FileSystem struct {
	basePath string
}

// NewFileSystem returns a new temporary filesystem.
// Creates the basePath if it not exists already.
// TODO: what if basePath already exists?
func NewFileSystem() *FileSystem {
	var basePath = fmt.Sprintf("/tmp/%s", StringWithCharset(8, "abcdefgh"))

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	return &FileSystem{basePath: basePath}
}

// BasePath returns the base path of the temp file system.
func (f *FileSystem) BasePath() string {
	return f.basePath
}

// Path returns the absolute file path for the given path.
func (f *FileSystem) Path(path string) string {
	return filepath.Join(f.basePath, path)
}

// CreateFromFileStructure creates the files and directories from the given
// fileStructure.
func (f *FileSystem) CreateFromFileStructure(fileStructure FileStructure) error {
	for _, path := range fileStructure {
		if path.IsDir() {
			err := f.MkdirAll(string(path))
			if err != nil {
				return err
			}
		} else {
			err := f.Create(string(path))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Cleanup removes all leftover files in the temp file system.
func (f *FileSystem) Cleanup() error {
	return os.RemoveAll(f.basePath)
}

// MkdirAll wraps os.MkdirAll().
func (f *FileSystem) MkdirAll(path string) error {
	return os.MkdirAll(f.Path(path), os.ModePerm)
}

// Create creates a new file. Also creates neccecary directories.
func (f *FileSystem) Create(path string) error {
	s := strings.Split(path, "/")
	if len(s) > 1 {
		sb := strings.Builder{}
		for _, dir := range s[:len(s)-1] {
			sb.WriteString(dir)
			f.MkdirAll(sb.String())
			sb.WriteByte('/')
		}
	}

	return ioutil.WriteFile(f.Path(path), []byte(""), os.ModePerm)
}

// TODO: look at this again
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
