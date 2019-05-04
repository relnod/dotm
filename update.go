package dotm

import (
	"context"
)

// UpdateOptions are the options used to update a dotfile profile.
type UpdateOptions struct {
	FromRemote bool
	ExecHooks  bool

	LinkOptions
}

// Update updates the symlinks for the given profile.
func Update(profile string, opts *UpdateOptions) error {
	c, err := LoadOrCreateConfig()
	if err != nil {
		return err
	}
	p, err := c.Profile(profile)
	if err != nil {
		return err
	}

	err = p.expandVars()
	if err != nil {
		return err
	}

	var hooks *Hooks
	if opts.ExecHooks {
		hooks, err = p.findHooks(&opts.TraversalOptions)
		if err != nil {
			return err
		}

		err = hooks.PreUpdate.Exec(context.Background())
		if err != nil {
			return err
		}
	}

	if opts.FromRemote {
		err = p.pullRemote()
		if err != nil {
			return err
		}
	}

	err = p.link(opts.LinkOptions)
	if err != nil {
		return err
	}

	if opts.ExecHooks {
		err = hooks.PostUpdate.Exec(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}
