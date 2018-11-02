package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
)

// New creates a set of new dotfile profiles.
func New(c *config.Config, names []string, configPath string) error {
	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		err = p.Create()
		if err != nil {
			return err
		}
	}

	err = c.WriteFile(configPath)
	if err != nil {
		return err
	}

	return nil
}
