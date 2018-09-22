package dotfiles

import (
	"io/ioutil"
	"path/filepath"

	"github.com/relnod/dotm/pkg/fileutil"
)

var defaultExcluded = []string{
	".git",
}

// Action defines an action, that can be run during the dotfile traversal.
type Action interface {
	// Run will get called for each file that gets traversed.
	Run(source, dest, name string) error
}

// Traverser is used to traverse the dotfiles structure.
type Traverser struct {
	excluded []string
}

// NewTraverser returns a new traverser.
func NewTraverser(excluded []string) *Traverser {
	return &Traverser{excluded: append(defaultExcluded, excluded...)}
}

// traverseVisitor implements the visitor.Interface.
type traverseVisitor struct {
	action Action
	source string
	dest   string
	name   string
}

// Visit calls the traversal action.
func (t traverseVisitor) Visit(dir, file string) {
	t.action.Run(
		filepath.Join(t.source, t.name, dir),
		filepath.Join(t.dest, dir),
		file,
	)
}

// Traverse traverses the dotfiles directory. Calling action.Run()
// for every file passed
// TODO: rethink arguments, maybe add Traverser struct
// TODO: finish implementation
func (t *Traverser) Traverse(source string, dest string, action Action) error {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		// TODO: wrap error
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		if t.isExcluded(f.Name()) {
			continue
		}

		t := traverseVisitor{
			action: action,
			source: source,
			dest:   dest,
			name:   f.Name(),
		}

		err := fileutil.RecTraverseDir(filepath.Join(source, f.Name()), "", t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Traverser) isExcluded(dir string) bool {
	for _, exclude := range t.excluded {
		if dir == exclude {
			return true
		}
	}

	return false
}
