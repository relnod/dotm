package dotfiles

import (
	"os/user"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Update updates the dotfiles for a given configuration.
func Update(c *config.Config, opts *UpdateOptions) error {
	var err error

	if opts == nil {
		opts = defaultUpdateOptions
	}

	err = c.Validate()
	if err != nil {
		return err
	}

	if opts.UpdateFromRemote {
		err = remote.Pull(c.FS, c.Remote, c.Path)
		if err != nil {
			return err
		}
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	err = Link(c.FS, c.Path, usr.HomeDir, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOptions is set of options for the update function.
type UpdateOptions struct {
	UpdateFromRemote bool
}

var defaultUpdateOptions = &UpdateOptions{
	UpdateFromRemote: false,
}
