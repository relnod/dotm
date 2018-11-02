package dotfiles

import (
	"errors"

	"github.com/relnod/fsa/testutil"

	"github.com/relnod/dotm/pkg/config"
)

// Errors
var (
	ErrPathExists = errors.New("path already exists")
)

// Install clones the dotfiles from a remote git repository to a local
// path and then installs them.
func Install(c *config.Config, names []string, configPath string, opts *InstallOptions) error {
	if opts == nil {
		opts = defaultInstallOptions
	}

	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		exists := testutil.FileExists(c.FS, p.Path)
		if exists {
			return ErrPathExists
		}

		err = p.CloneRemote()
		if err != nil {
			return err
		}

		err = p.Link(opts)
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

// InstallOptions is set of options for the install function. Implements the
// dotfiles.LinkOptions.
type InstallOptions struct {
	Force bool
	Dry   bool
}

// OptDry implementation
func (i *InstallOptions) OptDry() bool { return i.Dry }

// OptForce implementation
func (i *InstallOptions) OptForce() bool { return i.Force }

var defaultInstallOptions = &InstallOptions{
	Force: false,
	Dry:   false,
}
