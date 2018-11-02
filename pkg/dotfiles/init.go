package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
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
		err = p.Link(opts)
		if err != nil {
			return err
		}

		if p.Remote == "" {
			// Ignore error since remote detection is optional.
			remoteURL, _ := p.DetectRemote()
			p.Remote = remoteURL
		}
	}

	err = c.WriteFile(configPath)
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
