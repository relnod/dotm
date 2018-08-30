package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Install clones the dotfiles from a remote git repository to a local
// path and then installs them.
func Install(remoteURL string, path string) error {
	var err error

	c := &config.Config{
		Remote: remoteURL,
		Path:   path,
	}

	err = c.Validate()
	if err != nil {
		return err
	}

	// @TODO: check if dotfiles are already installed

	err = remote.Clone(c.Remote, c.Path)
	if err != nil {
		return err
	}

	err = Link(c.Path, "tmp/bla2", nil)

	err = config.WriteTomlFile("/tmp/bla3/.dotfiles.toml", c)
	if err != nil {
		return err
	}

	return nil
}
