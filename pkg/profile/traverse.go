package profile

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"

	"github.com/relnod/dotm/pkg/fileutil"
)

// Errors, that can occur during traversal
const (
	ErrReadDirPath = "failed to read dir for profile path"
)

// defaultExcluded specifies a list of directories, that will always be
// excluded.
var defaultExcluded = []string{".git"}

// Action defines an action, that can be run during the dotfile traversal.
type Action interface {
	// Run will get called for each file that gets traversed.
	Run(source, dest, name string) error
}

// Traverse traverses the dotfiles directory for a profile. Calls action.Run()
// for each passed file.
func Traverse(fs fsa.FileSystem, p *Profile, action Action) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	dest := usr.HomeDir

	traverseTopLevelDirs(fs, p, func(dir string) error {
		tv := traverseVisitor{
			action: action,
			source: p.Path,
			dest:   dest,
			name:   dir,
		}

		return fileutil.RecTraverseDir(fs, filepath.Join(p.Path, dir), "", tv)
	})

	return nil
}

func traverseTopLevelDirs(fs fsa.FileSystem, p *Profile, fn func(string) error) error {
	files, err := fsutil.ReadDir(fs, p.Path)
	if err != nil {
		return errors.Wrap(err, ErrReadDirPath)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if !isIncluded(p.Includes, f.Name()) {
			continue
		}
		if isExcluded(p.Excludes, f.Name()) {
			continue
		}

		err := fn(f.Name())
		if err != nil {
			return err
		}
	}

	return nil
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

// isExcluded checks if the directory should be excluded.
// Also excludes all directories prefixed with a "_".
func isExcluded(excludes []string, dir string) bool {
	if strings.HasPrefix(dir, "_") {
		return true
	}
	excludes = append(excludes, defaultExcluded...)
	for _, exclude := range excludes {
		if dir == exclude {
			return true
		}
	}
	return false
}

// isIncluded checks if the directory should be included.
func isIncluded(includes []string, dir string) bool {
	// If no includes are, all are included.
	if includes == nil {
		return true
	}
	for _, include := range includes {
		if dir == include {
			return true
		}
	}
	return false
}
