package dotfiles

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"
	"github.com/relnod/fsa/testutil"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/profile"
)

// Errors
var (
	ErrOpenHooksFile   = "failed to open hooks file '%s'"
	ErrDecodeHooksFile = "failed to decode hooks file '%s'"
	ErrExecuteHook     = "failed to execute hook '%s'"
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
		hooks, err := getHooks(c.FS, p.Path)
		if err != nil {
			return err
		}
		preUpdate := p.PreUpdate
		postUpdate := p.PostUpdate
		for _, h := range hooks {
			preUpdate = append(preUpdate, h.PreUpdate...)
			postUpdate = append(postUpdate, h.PostUpdate...)
		}

		err = executeHook(preUpdate)
		if err != nil {
			return err
		}

		if opts.UpdateFromRemote {
			err = profile.PullRemote(c.FS, p)
			if err != nil {
				return err
			}
		}

		err = profile.Link(c.FS, p, opts)
		if err != nil {
			return err
		}

		err = executeHook(postUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateOptions is set of options for the update function.
type UpdateOptions struct {
	UpdateFromRemote bool
	Force            bool
	Dry              bool
}

// OptDry implementation
func (i *UpdateOptions) OptDry() bool { return i.Dry }

// OptForce implementation
func (i *UpdateOptions) OptForce() bool { return i.Force }

var defaultUpdateOptions = &UpdateOptions{
	UpdateFromRemote: false,
	Force:            false,
	Dry:              false,
}

func executeHook(cmds []string) error {
	for _, cmdRaw := range cmds {
		args := strings.Split(os.ExpandEnv(cmdRaw), " ")
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return errors.Wrapf(err, ErrExecuteHook, cmdRaw)
		}
	}
	return nil
}

type hook struct {
	PreUpdate  []string `toml:"pre_update"`
	PostUpdate []string `toml:"post_update"`
}

func getHooks(fs fsa.FileSystem, dir string) ([]*hook, error) {
	var hooks []*hook

	h, err := findHook(fs, dir)
	if err != nil {
		return nil, err
	}
	if h != nil {
		hooks = append(hooks, h)
	}

	files, err := fsutil.ReadDir(fs, dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		h, err := findHook(fs, filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		if h != nil {
			hooks = append(hooks, h)
		}
	}
	return hooks, nil
}

func findHook(fs fsa.FileSystem, dir string) (*hook, error) {

	path := filepath.Join(dir, "hooks.toml")
	if !testutil.FileExists(fs, path) {
		return nil, nil
	}

	data, err := fsutil.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrapf(err, ErrOpenHooksFile, path)
	}

	var h hook
	_, err = toml.Decode(string(data), &h)
	if err != nil {
		return nil, errors.Wrapf(err, ErrDecodeHooksFile, path)
	}

	return &h, nil
}
