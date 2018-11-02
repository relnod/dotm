package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/hook"
)

// Update updates the dotfiles for a given configuration.
func Update(c *config.Config, names []string, opts *UpdateOptions) error {
	if opts == nil {
		opts = defaultUpdateOptions
	}

	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		var hooks *hook.Hooks
		if opts.Hooks {
			hooks, err = p.FindHooks()
			if err != nil {
				return err
			}

			err = hooks.PreUpdate.Execute()
			if err != nil {
				return err
			}
		}

		if opts.UpdateFromRemote && p.Remote != "" {
			err = p.PullRemote()
			if err != nil {
				return err
			}
		}

		err = p.Link(opts)
		if err != nil {
			return err
		}

		if opts.Hooks {
			err = hooks.PostUpdate.Execute()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateOptions is set of options for the update function.
type UpdateOptions struct {
	UpdateFromRemote bool
	Force            bool
	Dry              bool
	Hooks            bool
}

// OptDry implementation
func (i *UpdateOptions) OptDry() bool { return i.Dry }

// OptForce implementation
func (i *UpdateOptions) OptForce() bool { return i.Force }

var defaultUpdateOptions = &UpdateOptions{
	UpdateFromRemote: false,
	Force:            false,
	Dry:              false,
	Hooks:            true,
}
