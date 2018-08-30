package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Update updates the dotfiles for a given configuration.
func Update(c *config.Config) error {
	var err error

	err = c.Validate()
	if err != nil {
		return err
	}

	err = remote.Pull(c.Remote, c.Path)
	if err != nil {
		return err
	}

	err = Link(c.Path, "/tmp/bla2", nil)
	if err != nil {
		return err
	}

	return nil
}
