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

// Link implements the action.Interface for a link action.
type Link struct {
	dry bool
}

// NewLinkAction returns a new link action. If dry is set to true a dry run
// will be performed.
func NewLinkAction(dry bool) *Link {
	return &Link{dry: dry}
}

// Run links a file from source to dest.
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

		// TODO: override option (force)
		// TODO: backup option
	}

	if !os.IsNotExist(err) {
		return ErrorReadingFileStats
	}

	return file.Link(sourceFile, destFile, l.dry)
}

// Unlink implements the action.Interface for an unlink action.
type Unlink struct {
	dry bool
}

// NewUnlinkAction returns a new unlink action. If dry is set to true a dry run
// will be performed.
func NewUnlinkAction(dry bool) *Unlink {
	return &Unlink{dry: dry}
}

// Run unlinks the file a dest.
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
