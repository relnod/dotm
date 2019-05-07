package dotm

import (
	"context"
)

// InstallOptions are the options used to install a dotfile profile.
type InstallOptions struct {
	LinkOptions
}

// Install calls InstallWithContext with the background context.
func Install(p *Profile, opts *InstallOptions) error {
	return InstallWithContext(context.Background(), p, opts)
}

// InstallWithContext installs a new dotfile profile by cloning the remote
// repository to the local path. The clone operation can be canceled with the
// passed context.
func InstallWithContext(ctx context.Context, p *Profile, opts *InstallOptions) error {
	c, err := LoadOrCreateConfig()
	if err != nil {
		return err
	}

	p, err = c.AddProfile(p)
	if err != nil {
		return err
	}

	err = p.cloneRemote(ctx)
	if err != nil {
		return err
	}

	err = p.link(opts.LinkOptions)
	if err != nil {
		return err
	}

	c.Write()

	return nil
}
