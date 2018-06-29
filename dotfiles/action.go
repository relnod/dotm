package dotfiles

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/relnod/dotm/internal/util/file"
)

// TODO: use viper for configs
var excludes = []string{".git"}
var dry = true
var verbose = true

// Errors, that can occur during an action.
var (
	ErrorCreatingDestination = errors.New("Could not create destination directory")
	ErrorReadingFileStats    = errors.New("Could not read file stats")
)

// Action defines an action.
// TODO: improve doc
type Action interface {
	// TODO: add doc
	Run(source, dest, name string) error
}

// TODO: add doc
type Link struct {
	dry bool
}

func NewLinkAction(dry bool) *Link {
	return &Link{dry: dry}
}

// TODO: add doc
// TODO: add test
func (l *Link) Run(source, dest, name string) error {
	fmt.Println(source, dest, name)
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return ErrorCreatingDestination
	}

	sourceFile := filepath.Join(source, name)
	destFile := filepath.Join(dest, name)

	f, err := os.Stat(destFile)
	if err == nil {
		if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}

		// todo: override option (force)
		// todo: backup option
	}

	if !os.IsNotExist(err) {
		return ErrorReadingFileStats
	}

	return file.Link(sourceFile, destFile, l.dry)
}

// TODO: add doc
type Unlink struct {
	dry bool
}

func NewUnlinkAction(dry bool) *Unlink {
	return &Unlink{dry: dry}
}

// TODO: add doc
// TODO: add test
func (u *Unlink) Run(source, dest, name string) error {
	f, err := os.Stat(filepath.Join(dest, name))
	if err != nil {
		return nil
	}

	if f.Mode()&os.ModeSymlink != os.ModeSymlink {
		return nil
	}

	return file.Unlink(filepath.Join(dest, name), u.dry)
}
