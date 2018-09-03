package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
)

// Init initializes a set of dotfiles.
func Init(c *config.Config) error {
	var err error

	// @TODO: check if dotfiles are already installed

	err = Link(c.Path, "/tmp/bla2", nil)

	err = config.WriteTomlFile("/tmp/bla3/.dotfiles.toml", c)
	if err != nil {
		return err
	}

	return nil
}
