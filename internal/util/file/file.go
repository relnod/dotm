package file

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Visitor defines a visitor.
type Visitor interface {
	// TODO: add doc
	Visit(string, string)
}

type DefaultVisitor struct {
	fn func(string, string)
}

func NewDefaultVisitor(fn func(string, string)) *DefaultVisitor {
	return &DefaultVisitor{fn: fn}
}

func (d *DefaultVisitor) Visit(dir string, file string) {
	d.fn(dir, file)
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
// TODO: add tests
// TODO: return error
func Link(from string, to string, dry bool) {
	if dry {
		log.Printf("Creating Symlink from %s to %s", from, to)
		return
	}

	err := os.Symlink(from, to)
	if err != nil {
		log.Fatal("Error creating symlink!", err)
	}
}

// Unlink removes a symbolik link from dest to source. When dry is true only
// perfomers a dry run.
// TODO: add tests
// TODO: return error
func Unlink(file string, dry bool) {
	if dry {
		log.Printf("Removing %s", file)
		return
	}

	err := os.Remove(file)
	if err != nil {
		log.Fatal("Error while removing symlink!")
	}
}
