package dotfiles

import (
	"log"
	"os"
	"path/filepath"

	"github.com/relnod/dotm/internal/util/file"
)

// TODO: use viper for configs
var excludes = []string{".git"}
var dry = true
var verbose = true

// Action defines an action.
// TODO: improve doc
type Action interface {
	// TODO: add doc
	Run(source, dest, name string)
}

// TODO: add doc
type Link struct{}

// TODO: add doc
// TODO: add test
func (l *Link) Run(source, dest, name string) {
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating destination directory!")
	}

	f, err := os.Stat(filepath.Join(dest, name))
	if err == nil {
		if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			return
		}

		// todo: override option (force)
		// todo: backup option
	}

	if !os.IsNotExist(err) {
		log.Fatal("Error reading file stats", err)
	}

	file.Link(filepath.Join(source, name), filepath.Join(dest, name), dry)
}

// TODO: add doc
type Unlink struct{}

// TODO: add doc
// TODO: add test
func (u *Unlink) Run(source, dest, name string) {
	f, err := os.Stat(filepath.Join(dest, name))
	if err != nil {
		return
	}

	if f.Mode()&os.ModeSymlink != os.ModeSymlink {
		return
	}

	file.Unlink(filepath.Join(dest, name), dry)
}
