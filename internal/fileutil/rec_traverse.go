package fileutil

import (
	"io/ioutil"
	"path/filepath"
)

// Visitor defines a file visitor.
type Visitor interface {
	// Visit gets called for each file. The path is a relative path from the
	// start of the traversal. The name is the file name beeing visited.
	Visit(path string, name string) error
}

// RecTraverseDir recursively traverses all directories starting at dir.
// Calls the visitor.Visit for each file it passes.
func RecTraverseDir(dir string, visitor Visitor) error {
	return recTraverseDir(dir, "", visitor)
}

func recTraverseDir(dir, dirRelativeFromStart string, visitor Visitor) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			err := recTraverseDir(
				filepath.Join(dir, f.Name()),
				filepath.Join(dirRelativeFromStart, f.Name()),
				visitor,
			)
			if err != nil {
				return err
			}
		} else {
			err := visitor.Visit(dirRelativeFromStart, f.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
