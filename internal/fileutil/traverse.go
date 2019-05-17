package fileutil

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Visitor defines a file visitor.
type Visitor interface {
	// Visit gets called for each file. The path is a relative path from the
	// start of the traversal. The name is the file name beeing visited.
	Visit(path, relativePath string) error
}

// RecTraverseDir recursively traverses all directories starting at dir.
// Calls the visitor.Visit for each file it passes.
//
// All directories or files prefixed with the ignorePrefix are ignored.
func RecTraverseDir(dir string, visitor Visitor, ignorePrefix string) error {
	return recTraverseDir(dir, "", visitor, ignorePrefix)
}

func recTraverseDir(dir, relativeDir string, visitor Visitor, ignorePrefix string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if ignorePrefix != "" && strings.HasPrefix(f.Name(), ignorePrefix) {
			continue
		}
		if f.IsDir() {
			err := recTraverseDir(
				filepath.Join(dir, f.Name()),
				filepath.Join(relativeDir, f.Name()),
				visitor,
				ignorePrefix,
			)
			if err != nil {
				return err
			}
		} else {
			err := visitor.Visit(filepath.Join(dir, f.Name()), filepath.Join(relativeDir, f.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
