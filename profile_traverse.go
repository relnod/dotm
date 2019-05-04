package dotm

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/relnod/dotm/internal/fileutil"
)

// TraversalOptions are used during the dotfile traversal.
type TraversalOptions struct {
	Includes []string
	Excludes []string
}

func (p *Profile) traverse(visitor fileutil.Visitor, opts *TraversalOptions) error {
	topLevelDirs, err := p.topLevelDirs(opts)
	if err != nil {
		return err
	}

	for _, d := range topLevelDirs {
		err := fileutil.RecTraverseDir(filepath.Join(p.expandedPath, d), visitor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) topLevelDirs(opts *TraversalOptions) ([]string, error) {
	files, err := ioutil.ReadDir(p.expandedPath)
	if err != nil {
		return nil, err
	}

	includes := p.Includes
	excludes := p.Excludes
	if opts != nil {
		includes = append(includes, opts.Includes...)
		excludes = append(excludes, opts.Excludes...)
	}

	dirs := []string{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if !isIncluded(f.Name(), includes) {
			continue
		}
		if isExcluded(f.Name(), excludes) {
			continue
		}

		dirs = append(dirs, f.Name())
	}

	return dirs, nil
}

// alwaysExcluded is a list of directories, that always get excluded.
var alwaysExcluded = []string{".git"}

// isExcluded checks if the dir should be excluded.
// Also excludes all directories prefixed with a "_".
func isExcluded(dir string, excludes []string) bool {
	if strings.HasPrefix(dir, "_") {
		return true
	}
	excludes = append(excludes, alwaysExcluded...)
	for _, exclude := range excludes {
		if dir == exclude {
			return true
		}
	}
	return false
}

// isIncluded checks if the directory should be included.
// When the includes list is empty, it returns true.
func isIncluded(dir string, includes []string) bool {
	if len(includes) == 0 {
		return true
	}
	for _, include := range includes {
		if dir == include {
			return true
		}
	}
	return false
}
