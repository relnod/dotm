package dotm

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
	_, err = c.Profile(p.Name)
	if err == nil {
		return ErrProfileExists
	}

	err = p.expandVars()
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

	c.Profiles[p.Name] = p
	c.Write()

	return nil
}
