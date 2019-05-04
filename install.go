package dotm

import (
	"errors"
	"os"
)

// InstallOptions are the options used to install a dotfile profile.
type InstallOptions struct {
	LinkOptions
}

// ErrProfileExists indicates, that the profile already exists.
var ErrProfileExists = errors.New("profile already exists")

// ErrProfilePathExists indicates, that the profile path already exists.
var ErrProfilePathExists = errors.New("profile path already exists")

// Install installs a new dotfile profile by cloning the remote repository to
// the local path.
func Install(p *Profile, opts *InstallOptions) error {
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

	err = p.cloneRemote()
	if err != nil {
		return err
	}

	err = p.link(opts.LinkOptions)
	if err != nil {
		return err
	}

	c.Profiles[p.Name] = p
	c.Write()

	return nil
}
