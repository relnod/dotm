package dotfiles

import (
	"errors"

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
func Install(c *config.Config, names []string, configPath string) error {
	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		exists := testutil.FileExists(c.FS, p.Path)
		if exists {
			return ErrPathExists
		}

		err = remote.CloneProfile(c.FS, p)
		if err != nil {
			return err
		}

		err = LinkProfile(c.FS, p)
		if err != nil {
			return err
		}
	}

	err = config.WriteFile(c.FS, configPath, c)
	if err != nil {
		return err
	}

	return nil
}
