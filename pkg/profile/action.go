package profile

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/relnod/fsa"
	"github.com/relnod/fsa/testutil"

	"github.com/relnod/dotm/pkg/fileutil"
)

// Errors, that can occur during an action.
var (
	ErrCreatingDestination = errors.New("Could not create destination directory")
	ErrReadingFileStats    = errors.New("Could not read file stats")
)

// Link recursively links all files from one profile.
func (p *Profile) Link(opts LinkOptions) error {
	err := p.Traverse(NewLinkAction(p.fs, opts))
	if err != nil {
		return err
	}
	return nil
}

// Unlink recursively removes the symlinks for one profile.
func (p *Profile) Unlink(opts UnlinkOptions) error {
	err := p.Traverse(NewUnlinkAction(p.fs, opts))
	if err != nil {
		return err
	}
	return nil
}

// ActionOptions defines the generic options for an action.
type ActionOptions interface {
	// Dry specifies if the action should perform a dry run.
	OptDry() bool
}

// LinkOptions defines how an option for a link action should look.
type LinkOptions interface {
	ActionOptions

	// Force specifies if existing files should be overwritten by a symlink.
	OptForce() bool
}

type defaultLinkOptions struct{}

func (d *defaultLinkOptions) OptDry() bool   { return false }
func (d *defaultLinkOptions) OptForce() bool { return false }

// LinkAction implements the action.Interface for a link action.
type LinkAction struct {
	fs   fsa.FileSystem
	opts LinkOptions
}

// NewLinkAction returns a new link action.
func NewLinkAction(fs fsa.FileSystem, opts LinkOptions) *LinkAction {
	if opts == nil {
		opts = &defaultLinkOptions{}
	}
	return &LinkAction{fs: fs, opts: opts}
}

// Run links a file from source to dest.
func (l *LinkAction) Run(source, dest, name string) error {
	err := l.fs.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return ErrCreatingDestination
	}

	sourceFile := filepath.Join(source, name)
	destFile := filepath.Join(dest, name)

	if testutil.FileExists(l.fs, destFile) {
		if !l.opts.OptForce() {
			return nil
		}
		if testutil.IsSymlink(l.fs, destFile) {
			fileutil.Unlink(l.fs, destFile, l.opts.OptDry())

		} else {
			fileutil.Backup(l.fs, destFile, l.opts.OptDry())

		}
	}

	return fileutil.Link(l.fs, sourceFile, destFile, l.opts.OptDry())
}

// UnlinkOptions defines how an option for an unlink action should look.
type UnlinkOptions interface {
	ActionOptions
}

type defaultUnlinkOptions struct{}

func (d *defaultUnlinkOptions) OptDry() bool { return false }

// UnlinkAction implements dotfiles.Action with an unlink action.
type UnlinkAction struct {
	fs   fsa.FileSystem
	opts UnlinkOptions
}

// NewUnlinkAction returns a new unlink action.
func NewUnlinkAction(fs fsa.FileSystem, opts UnlinkOptions) *UnlinkAction {
	if opts == nil {
		opts = &defaultUnlinkOptions{}
	}
	return &UnlinkAction{fs: fs, opts: opts}
}

// Run unlinks the file at dest.
func (u *UnlinkAction) Run(source, dest, name string) error {
	path := filepath.Join(dest, name)
	fi, err := u.fs.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
		return nil
	}

	err = fileutil.Unlink(u.fs, path, u.opts.OptDry())
	if err != nil {
		return err
	}
	return fileutil.RestoreBackup(u.fs, path, u.opts.OptDry())
}
