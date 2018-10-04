package dotfiles

import (
	"path/filepath"

	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"

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
	fs       fsa.FileSystem
	excludes []string
	includes []string
}

// NewTraverser returns a new traverser.
func NewTraverser(fs fsa.FileSystem, excludes []string, includes []string) *Traverser {
	return &Traverser{
		fs:       fs,
		excludes: append(defaultExcluded, excludes...),
		includes: includes,
	}
}

// traverseVisitor implements fileutil.Visitor.
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
	files, err := fsutil.ReadDir(t.fs, source)
	if err != nil {
		// TODO: wrap error
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		if !t.isIncluded(f.Name()) {
			continue
		}

		if t.isExcluded(f.Name()) {
			continue
		}

		tv := traverseVisitor{
			action: action,
			source: source,
			dest:   dest,
			name:   f.Name(),
		}

		err := fileutil.RecTraverseDir(t.fs, filepath.Join(source, f.Name()), "", tv)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Traverser) isExcluded(dir string) bool {
	for _, exclude := range t.excludes {
		if dir == exclude {
			return true
		}
	}
	return false
}

func (t *Traverser) isIncluded(dir string) bool {
	// If no includes are, all are included.
	if t.includes == nil {
		return true
	}

	for _, include := range t.includes {
		if dir == include {
			return true
		}
	}
	return false
}
