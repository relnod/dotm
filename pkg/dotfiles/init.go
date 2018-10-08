package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/profile"
)

// Init initializes a set of dotfiles.
func Init(c *config.Config, names []string, configPath string, opts *InitOptions) error {
	if opts == nil {
		opts = defaultInitOptions
	}

	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	// TODO: check if dotfiles are already installed

	for _, p := range profiles {
		err = profile.Link(c.FS, p, opts)
		if err != nil {
			return err
		}

		if p.Remote == "" {
			// Ignore error since remote detection is optional.
			remoteURL, _ := profile.DetectRemote(c.FS, p)
			p.Remote = remoteURL
		}
	}

	err = config.WriteFile(c.FS, configPath, c)
	if err != nil {
		return err
	}

	return nil
}

// InitOptions is set of options for the init function. Implements the
// dotfiles.LinkOptions.
type InitOptions struct {
	Force bool
	Dry   bool
}

// OptDry implementation
func (i *InitOptions) OptDry() bool { return i.Dry }

// OptForce implementation
func (i *InitOptions) OptForce() bool { return i.Force }

var defaultInitOptions = &InitOptions{
	Force: false,
	Dry:   false,
}
