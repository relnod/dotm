package profile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/relnod/fsa"

	"github.com/relnod/dotm/pkg/hook"
)

// Profile defines one set of dotfiles.
type Profile struct {
	Remote   string   `toml:"remote"`
	Path     string   `toml:"path"`
	Includes []string `toml:"includes"`
	Excludes []string `toml:"excludes"`
	hook.Hooks

	fs fsa.FileSystem
}

// SetFS sets the file system.
func (p *Profile) SetFS(fs fsa.FileSystem) {
	p.fs = fs
}

// ExpandEnvs expands several variables with environment variables.
func (p *Profile) ExpandEnvs(name string) (err error) {
	p.Remote = os.ExpandEnv(p.Remote)
	p.Path, err = expandPath(p.Path, name)
	return err
}

func expandPath(path, name string) (string, error) {
	path = strings.Replace(path, "<PROFILE>", name, 1)

	path = os.ExpandEnv(path)

	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}
