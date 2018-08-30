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

func Link(from, to string, opts *LinkOptions) error {
	if opts == nil {
		opts = defaultLinkOptions
	}

	t := NewTraverser(nil)
	err := t.Traverse(from, to, NewLinkAction(opts.Dry))
	if err != nil {
		return err
	}

	return nil
}

type LinkOptions struct {
	Dry bool
}

var defaultLinkOptions = &LinkOptions{
	Dry: false,
}

// LinkAction implements the action.Interface for a link action.
type LinkAction struct {
	dry bool
}

// NewLinkAction returns a new link action. If dry is set to true a dry run
// will be performed.
func NewLinkAction(dry bool) *LinkAction {
	return &LinkAction{dry: dry}
}

// Run links a file from source to dest.
func (l *LinkAction) Run(source, dest, name string) error {
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

// UnlinkAction implements the action.Interface for an unlink action.
type UnlinkAction struct {
	dry bool
}

// NewUnlinkAction returns a new unlink action. If dry is set to true a dry run
// will be performed.
func NewUnlinkAction(dry bool) *UnlinkAction {
	return &UnlinkAction{dry: dry}
}

// Run unlinks the file a dest.
func (u *UnlinkAction) Run(source, dest, name string) error {
	f, err := os.Stat(filepath.Join(dest, name))
	if err != nil {
		return nil
	}

	if f.Mode()&os.ModeSymlink != os.ModeSymlink {
		return nil
	}

	return file.Unlink(filepath.Join(dest, name), u.dry)
}
