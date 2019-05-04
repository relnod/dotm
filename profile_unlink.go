package dotm

import (
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/relnod/dotm/internal/file"
)

func (p *Profile) unlink(dry bool) error {
	err := p.traverse(&unlinker{
		path: os.Getenv("HOME"),
		dry:  dry,
	}, nil)
	if err != nil {
		return xerrors.Errorf("unlink: ", err)
	}
	return nil
}

// unlinker implements fileutil.Visitor
type unlinker struct {
	path string

	dry bool
}

func (u *unlinker) Visit(path, name string) error {
	filepath := filepath.Join(u.path, path, name)

	// Check if the file file exists.
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}
	if !file.IsSymlink(filepath) {
		return nil
	}

	err := file.Unlink(filepath, u.dry)
	if err != nil {
		return err
	}

	return file.RestoreBackup(filepath, u.dry)
}
