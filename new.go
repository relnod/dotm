package dotm

import (
	"os"
)

// New creates a new dotfile profile.
func New(p *Profile) error {
	c, err := LoadOrCreateConfig()
	if err != nil {
		return err
	}
	if _, ok := c.Profiles[p.Name]; ok {
		return ErrProfileExists
	}

	err = p.expandVars()
	if err != nil {
		return err
	}

	if _, err := os.Stat(p.expandedPath); err == nil {
		return ErrProfilePathExists
	}

	err = p.create()
	if err != nil {
		return err
	}

	c.Profiles[p.Name] = p
	c.Write()

	return nil
}
