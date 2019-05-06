package dotm

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

	err = p.expandEnv()
	if err != nil {
		return err
	}

	err = p.unlink(opts.Dry)
	if err != nil {
		return err
	}

	c.Profiles[p.Name] = p
	c.Write()

	return nil
}
