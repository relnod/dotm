package profile

import (
	"path/filepath"

	"github.com/relnod/fsa"
	"github.com/relnod/fsa/testutil"

	"github.com/relnod/dotm/pkg/hook"
)

// FindHooks finds all hooks of a given profile.
// Hooks can be found at:
// - ~/.dotfiles/dotm.toml
// - ~/.dotfiles/<profile>/hooks.toml
// - ~/.dotfiles/<profile>/<top-level-dir>/hooks.toml
func (p *Profile) FindHooks() (*hook.Hooks, error) {
	var hooks []*hook.Hooks

	hooks = append(hooks, &p.Hooks)

	h, err := findHook(p.fs, p.Path)
	if err != nil {
		return nil, err
	}
	if h != nil {
		hooks = append(hooks, h)
	}

	err = traverseTopLevelDirs(p.fs, p, func(dir string) error {
		h, err := findHook(p.fs, filepath.Join(p.Path, dir))
		if err != nil {
			return err
		}
		if h != nil {
			hooks = append(hooks, h)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return hook.Merge(hooks...), nil
}

func findHook(fs fsa.FileSystem, dir string) (*hook.Hooks, error) {
	path := filepath.Join(dir, "hooks.toml")
	if !testutil.FileExists(fs, path) {
		return nil, nil
	}

	return hook.ReadFromFile(fs, path)
}