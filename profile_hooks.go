package dotm

import (
	"os"
	"path/filepath"
)

// FindHooks finds all hooks of a given profile.
// Hooks can be found at:
// - ~/.dotfiles/dotm.toml
// - ~/.dotfiles/<profile>/hooks.toml
// - ~/.dotfiles/<profile>/<top-level-dir>/hooks.toml
func (p *Profile) findHooks(opts *TraversalOptions) (*Hooks, error) {
	var hooks []*Hooks

	hooks = append(hooks, &p.Hooks)

	h, err := findHook(p.expandedPath)
	if err != nil {
		return nil, err
	}
	if h != nil {
		hooks = append(hooks, h)
	}

	topLevelDirs, err := p.topLevelDirs(opts)
	if err != nil {
		return nil, err
	}
	for _, dir := range topLevelDirs {
		h, err := findHook(filepath.Join(p.expandedPath, dir))
		if err != nil {
			return nil, err
		}
		if h != nil {
			hooks = append(hooks, h)
		}
	}

	return mergeHooks(hooks...), nil
}

const hooksFileName = "hooks.toml"

func findHook(dir string) (*Hooks, error) {
	filepath := filepath.Join(dir, hooksFileName)
	if _, err := os.Stat(filepath); err != nil {
		return nil, nil
	}

	return LoadHooksFromFile(filepath)
}
