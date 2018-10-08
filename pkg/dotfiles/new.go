package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/profile"
)

// New creates a set of new dotfile profiles.
func New(c *config.Config, names []string, configPath string) error {
	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		err = profile.Create(c.FS, p)
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
