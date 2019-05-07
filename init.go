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

	p.sanitize()
	err = p.expandEnv()
	if err != nil {
		return err
	}
	err = c.AddProfileFromExistingPath(p)
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
