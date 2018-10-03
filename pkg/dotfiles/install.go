package dotfiles

import (
	"errors"
	"os/user"

	"github.com/relnod/fsa/testutil"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Errors
var (
	ErrPathExists = errors.New("path already exists")
)

// Install clones the dotfiles from a remote git repository to a local
// path and then installs them.
func Install(c *config.Config, configPath string) error {
	var err error

	err = c.Validate()
	if err != nil {
		return err
	}

	exists := testutil.FileExists(c.FS, c.Path)
	if exists {
		return ErrPathExists
	}

	err = remote.Clone(c.FS, c.Remote, c.Path)
	if err != nil {
		return err
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	err = Link(c.FS, c.Path, usr.HomeDir, nil)

	err = config.WriteFile(c.FS, configPath, c)
	if err != nil {
		return err
	}

	return nil
}
