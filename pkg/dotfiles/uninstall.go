package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
)

// Uninstall uninstalles the dotfiles.
func Uninstall(c *config.Config, names []string, opts *UninstallOptions) error {
	if opts == nil {
		opts = defaultUninstallOptions
	}

	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		err = p.Unlink(opts)
		if err != nil {
			return err
		}
		// err = os.RemoveAll(c.Path)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

// UninstallOptions is set of options for the uninstall function. Implements the
// dotfiles.UnlinkOptions.
type UninstallOptions struct {
	Force bool
	Dry   bool
}

// OptDry implementation
func (i *UninstallOptions) OptDry() bool { return i.Dry }

var defaultUninstallOptions = &UninstallOptions{
	Force: false,
	Dry:   false,
}
