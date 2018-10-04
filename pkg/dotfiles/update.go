package dotfiles

import (
	"os/exec"
	"strings"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Update updates the dotfiles for a given configuration.
func Update(c *config.Config, names []string, opts *UpdateOptions) error {
	if opts == nil {
		opts = defaultUpdateOptions
	}

	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		err = executeHook(p.PreUpdate)
		if err != nil {
			return err
		}

		if opts.UpdateFromRemote {
			err = remote.PullProfile(c.FS, p)
			if err != nil {
				return err
			}
		}

		err = LinkProfile(c.FS, p)
		if err != nil {
			return err
		}

		err = executeHook(p.PostUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateOptions is set of options for the update function.
type UpdateOptions struct {
	UpdateFromRemote bool
}

var defaultUpdateOptions = &UpdateOptions{
	UpdateFromRemote: false,
}

func executeHook(cmds []string) error {
	for _, cmdRaw := range cmds {
		args := strings.Split(cmdRaw, " ")
		cmd := exec.Command(args[0], args[1:]...)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
