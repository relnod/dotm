package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Init initializes a set of dotfiles.
func Init(c *config.Config, names []string, configPath string) error {
	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	// TODO: check if dotfiles are already installed

	for _, profile := range profiles {
		err = LinkProfile(c.FS, profile)
		if err != nil {
			return err
		}

		if profile.Remote == "" {
			// Ignore error since remote detection is optional.
			remoteURL, _ := remote.Detect(profile.Path)
			profile.Remote = remoteURL
		}
	}

	err = config.WriteFile(c.FS, configPath, c)
	if err != nil {
		return err
	}

	return nil
}
