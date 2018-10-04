package dotfiles

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"
	"github.com/relnod/fsa/testutil"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
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
			err = remote.PullProfile(c.FS, p)
			if err != nil {
				return err
			}
		}

		err = LinkProfile(c.FS, p)
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
			return errors.Wrapf(err, ErrExecuteHook, cmdRaw)
		}
	}
	return nil
}

type hook struct {
	PreUpdate  []string `toml:"pre_update"`
	PostUpdate []string `toml:"post_update"`
}

func getHooks(fs fsa.FileSystem, dir string) ([]hook, error) {
	files, err := fsutil.ReadDir(fs, dir)
	if err != nil {
		return nil, err
	}

	var hooks []hook
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		path := filepath.Join(dir, file.Name(), "hooks.toml")
		if !testutil.FileExists(fs, path) {
			continue
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
		hooks = append(hooks, h)
	}
	return hooks, nil
}
