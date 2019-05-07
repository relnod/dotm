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

// Update calls UpdateWithContext with the background context.
func Update(profile string, opts *UpdateOptions) error {
	return UpdateWithContext(context.Background(), profile, opts)
}

// UpdateWithContext updates the symlinks for the given profile. When
// opts.FromRemote is set to true it first pull updates from the remote
// repository. This operation can be canceled with the passed context.
// When opts.ExecHooks is passed, pre and post update hooks get executed.
func UpdateWithContext(ctx context.Context, profile string, opts *UpdateOptions) error {
	c, err := LoadConfig()
	if err != nil {
		return err
	}
	p, err := c.Profile(profile)
	if err != nil {
		return err
	}

	var hooks *Hooks
	if opts.ExecHooks {
		hooks, err = p.findHooks(&opts.TraversalOptions)
		if err != nil {
			return err
		}

		err = hooks.PreUpdate.Exec(ctx)
		if err != nil {
			return err
		}
	}

	if opts.FromRemote {
		err = p.pullRemote(ctx)
		if err != nil {
			return err
		}
	}

	err = p.link(opts.LinkOptions)
	if err != nil {
		return err
	}

	if opts.ExecHooks {
		err = hooks.PostUpdate.Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
