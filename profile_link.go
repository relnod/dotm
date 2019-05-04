package dotm

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/relnod/dotm/internal/file"
)

// LinkOptions are the options used during the symlink creation.
type LinkOptions struct {
	Force bool
	Dry   bool
	TraversalOptions
}

// link links all files to the home directory.
func (p *Profile) link(opts LinkOptions) error {
	err := p.traverse(&linker{
		source: p.expandedPath,
		dest:   os.Getenv("HOME"),
		force:  opts.Force,
		dry:    opts.Dry,
	}, &opts.TraversalOptions)
	if err != nil {
		return xerrors.Errorf("link: ", err)
	}
	return nil
}

// linker implements fileutil.Visitor
type linker struct {
	// source is the path from where to link from
	source string

	// dest is the path where the files get linked to
	dest string

	force bool
	dry   bool
}

//ErrCreatingDestination indicates failure during the creation of the
//destination dir
var ErrCreatingDestination = errors.New("failed to created destination dir")

func (l *linker) Visit(path, name string) error {
	err := os.MkdirAll(filepath.Join(l.dest, path), os.ModePerm)
	if err != nil {
		return ErrCreatingDestination
	}

	var (
		sourceFile = filepath.Join(l.source, path, name)
		destFile   = filepath.Join(l.dest, path, name)
	)

	// Check if the destination file already exists.
	if _, err := os.Stat(destFile); err == nil {
		if !l.force {
			return nil
		}
		if file.IsSymlink(destFile) {
			err = file.Unlink(destFile, l.dry)
			if err != nil {
				return err
			}
		} else {
			err = file.Backup(destFile, l.dry)
			if err != nil {
				return err
			}
		}
	}

	return file.Link(sourceFile, destFile, l.dry)
}
