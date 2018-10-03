package dotfiles

import (
	"os/user"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Init initializes a set of dotfiles.
func Init(c *config.Config, configPath string) error {
	var err error

	// TODO: check if dotfiles are already installed

	usr, err := user.Current()
	if err != nil {
		return err
	}

	err = Link(c.FS, c.Path, usr.HomeDir, nil)
	if err != nil {
		return err
	}

	if c.Remote == "" {
		// Ignore error since remote detection is optional.
		remoteURL, _ := remote.Detect(c.Path)
		c.Remote = remoteURL
	}

	err = config.WriteFile(c.FS, configPath, c)
	if err != nil {
		return err
	}

	return nil
}
