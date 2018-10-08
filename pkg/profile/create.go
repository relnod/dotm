package profile

import (
	"os"

	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	git "gopkg.in/src-d/go-git.v4"
)

// Errors
const (
	ErrInitRepo = "failed to initialize git repo"
)

// Create creates the path of the profile.
func Create(fs fsa.FileSystem, p *Profile) error {
	_, path := sanitizePaths(fs, p.Remote, p.Path)

	err := fs.MkdirAll(path, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrMkdirDestination)
	}

	_, err = git.PlainInit(path, false)
	if err != nil {
		return errors.Wrap(err, ErrInitRepo)
	}
	return nil
}
