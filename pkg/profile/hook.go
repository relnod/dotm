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
func FindHooks(fs fsa.FileSystem, p *Profile) (*hook.Hooks, error) {
	var hooks []*hook.Hooks

	hooks = append(hooks, &p.Hooks)

	h, err := findHook(fs, p.Path)
	if err != nil {
		return nil, err
	}
	if h != nil {
		hooks = append(hooks, h)
	}

	traverseTopLevelDirs(fs, p, func(dir string) error {
		h, err := findHook(fs, filepath.Join(p.Path, dir))
		if err != nil {
			return err
		}
		if h != nil {
			hooks = append(hooks, h)
		}
		return nil
	})

	return hook.Merge(hooks...), nil
}

func findHook(fs fsa.FileSystem, dir string) (*hook.Hooks, error) {
	path := filepath.Join(dir, "hooks.toml")
	if !testutil.FileExists(fs, path) {
		return nil, nil
	}

	return hook.ReadFromFile(fs, path)
}
