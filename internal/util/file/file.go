package file

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Errors
var (
	ErrCreatingSymlink = errors.New("Failed to create Symlink")
)

// Visitor defines a visitor.
type Visitor interface {
	// Visit gets called for each file being traversed.
	Visit(dir string, file string)
}

// RecTraverseDir recursively traverses all directories starting at dir.
// Calls the visitor for each file it passes.
// TODO: what is relDir for?
func RecTraverseDir(dir string, relDir string, visitor Visitor) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			err := RecTraverseDir(filepath.Join(dir, f.Name()), filepath.Join(relDir, f.Name()), visitor)
			if err != nil {
				return err
			}
		} else {
			visitor.Visit(relDir, f.Name())
		}
	}

	return nil
}

// Link creates a symbolik link from dest to source. When dry is true only
// perfomers a dry run.
func Link(from string, to string, dry bool) error {
	if dry {
		log.Printf("Creating Symlink from %s to %s", from, to)
		return nil
	}

	err := os.Symlink(from, to)
	if err != nil {
		return ErrCreatingSymlink
	}

	return nil
}

// Unlink removes a symbolik link from dest to source. When dry is true only
// perfomers a dry run.
// TODO: add tests
func Unlink(file string, dry bool) error {
	if dry {
		log.Printf("Removing %s", file)
		return nil
	}

	err := os.Remove(file)
	if err != nil {
		return err
	}

	return nil
}

// Exists checks if a file at the given path exists.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
