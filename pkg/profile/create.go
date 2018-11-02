package profile

import (
	"os"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
)

// Errors
const (
	ErrInitRepo = "failed to initialize git repo"
)

// Create creates the path of the profile.
func (p *Profile) Create() error {
	_, path := sanitizePaths(p.fs, p.Remote, p.Path)

	err := p.fs.MkdirAll(path, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrMkdirDestination)
	}

	_, err = git.PlainInit(path, false)
	if err != nil {
		return errors.Wrap(err, ErrInitRepo)
	}
	return nil
}
