package dotfiles

import (
	"io/ioutil"
	"path/filepath"

	"github.com/relnod/dotm/internal/util/file"
)

// Traverse traverses the dotfiles directory. Calling action.Run()
// for every file passed
// TODO: rethink arguments, maybe add Traverser struct
// TODO: finish implementation
func Traverse(source string, dest string, action Action) error {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		//TODO: wrap error
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		if isExcluded(f.Name()) {
			continue
		}

		v := file.NewDefaultVisitor(func(dir, file string) {
			action.Run(
				filepath.Join(source, f.Name(), dir),
				filepath.Join(dest, dir),
				file,
			)
		})

		err := file.RecTraverseDir(filepath.Join(source, f.Name()), "", v)
		if err != nil {
			return err
		}
	}

	return nil
}
