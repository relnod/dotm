package dotm

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
