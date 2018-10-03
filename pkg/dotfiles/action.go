package dotfiles

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/relnod/fsa"

	"github.com/relnod/dotm/pkg/fileutil"
)

// Errors, that can occur during an action.
var (
	ErrCreatingDestination = errors.New("Could not create destination directory")
	ErrReadingFileStats    = errors.New("Could not read file stats")
)

// Link recursively links files from one path to another.
func Link(fs fsa.FileSystem, from, to string, opts *ActionOptions) error {
	if opts == nil {
		opts = defaultActionOptions
	}

	t := NewTraverser(fs, nil)
	err := t.Traverse(from, to, NewLinkAction(fs, opts.Dry))
	if err != nil {
		return err
	}

	return nil
}

// UnLink recursively removes the symlinks.
func UnLink(fs fsa.FileSystem, from, to string, opts *ActionOptions) error {
	if opts == nil {
		opts = defaultActionOptions
	}

	t := NewTraverser(fs, nil)
	err := t.Traverse(from, to, NewUnlinkAction(fs, opts.Dry))
	if err != nil {
		return err
	}

	return nil
}

// ActionOptions defines a set options for an action
type ActionOptions struct {
	Dry      bool
	Verbose  bool
	Excludes []string
}

var defaultActionOptions = &ActionOptions{
	Dry:      false,
	Verbose:  false,
	Excludes: defaultExcluded,
}

// LinkAction implements the action.Interface for a link action.
type LinkAction struct {
	fs  fsa.FileSystem
	dry bool
}

// NewLinkAction returns a new link action. If dry is set to true a dry run
// will be performed.
func NewLinkAction(fs fsa.FileSystem, dry bool) *LinkAction {
	return &LinkAction{fs: fs, dry: dry}
}

// Run links a file from source to dest.
func (l *LinkAction) Run(source, dest, name string) error {
	err := l.fs.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return ErrCreatingDestination
	}

	sourceFile := filepath.Join(source, name)
	destFile := filepath.Join(dest, name)

	f, err := l.fs.Stat(destFile)
	if err == nil {
		if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}

		// TODO: override option (force)
		// TODO: backup option
	}
	if !os.IsNotExist(err) {
		return ErrReadingFileStats
	}

	return fileutil.Link(l.fs, sourceFile, destFile, l.dry)
}

// UnlinkAction implements the action.Interface for an unlink action.
type UnlinkAction struct {
	fs  fsa.FileSystem
	dry bool
}

// NewUnlinkAction returns a new unlink action. If dry is set to true a dry run
// will be performed.
func NewUnlinkAction(fs fsa.FileSystem, dry bool) *UnlinkAction {
	return &UnlinkAction{fs: fs, dry: dry}
}

// Run unlinks the file at dest.
func (u *UnlinkAction) Run(source, dest, name string) error {
	fi, err := u.fs.Lstat(filepath.Join(dest, name))
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
		return nil
	}

	return fileutil.Unlink(u.fs, filepath.Join(dest, name), u.dry)
}
