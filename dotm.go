package dotm

import "context"

// New creates a new dotfile profile.
func New(p *Profile) error {
	c, err := LoadOrCreateConfig()
	if err != nil {
		return err
	}

	p, err = c.AddProfile(p)
	if err != nil {
		return err
	}

	err = p.create()
	if err != nil {
		return err
	}

	c.Write()

	return nil
}

// InitOptions are the options used to initialize a dotfile profile.
type InitOptions struct {
	LinkOptions
}

// Init initializes a new dotfile profile from an existing local path.
func Init(p *Profile, opts *InitOptions) error {
	c, err := LoadOrCreateConfig()
	if err != nil {
		return err
	}

	p, err = c.AddProfileFromExistingPath(p)
	if err != nil {
		return err
	}

	err = p.link(opts.LinkOptions)
	if err != nil {
		return err
	}

	// Ignore error since remote detection is optional.
	remoteURL, _ := p.detectRemote()
	p.Remote = remoteURL

	c.Write()

	return nil
}

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

// Add adds the given file to the profile under the given top level dir.
// If the file already exists under $HOME/path, A backup is created and then
// copied to the profile.
func Add(profile, dir, path string) error {
	c, err := LoadConfig()
	if err != nil {
		return err
	}
	p, err := c.Profile(profile)
	if err != nil {
		return err
	}

	return p.addFile(dir, path)
}

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

// UpdateWithContext updates the symlinks for the given profile.
//
// When the given profile is empty all profiles get updated.
func UpdateWithContext(ctx context.Context, profile string, opts *UpdateOptions) error {
	c, err := LoadConfig()
	if err != nil {
		return err
	}

	// When the profile name is empty update all profiles.
	if profile == "" {
		for _, p := range c.Profiles {
			err = update(ctx, p, opts)
			if err != nil {
				return err
			}
		}
		return nil
	}

	p, err := c.Profile(profile)
	if err != nil {
		return err
	}
	return update(ctx, p, opts)
}

// update updates the symlinks for a given profile.
//
// When opts.FromRemote is set to true it first pull updates from the remote
// repository. This operation can be canceled with the passed context.
// When opts.ExecHooks is passed, pre and post update hooks get executed.
func update(ctx context.Context, p *Profile, opts *UpdateOptions) (err error) {
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

// UninstallOptions are the options used to uninstall a dotfile profile.
type UninstallOptions struct {
	Dry bool
}

// Uninstall unlinks a dotfile profile. If any backup files exist, they get
// restored.
func Uninstall(profile string, opts *UninstallOptions) error {
	c, err := LoadOrCreateConfig()
	if err != nil {
		return err
	}
	p, err := c.Profile(profile)
	if err != nil {
		return err
	}

	err = p.unlink(opts.Dry)
	if err != nil {
		return err
	}

	c.Write()

	return nil
}
