package dotm

import (
	"os"
	"path/filepath"
	"strings"
)

// Profile defines the data of a dotfile profile.
type Profile struct {
	Name     string   `toml:"-"`
	Path     string   `toml:"path"`
	Remote   string   `toml:"remote"`
	Includes []string `toml:"includes"`
	Excludes []string `toml:"excludes"`
	Hooks

	expandedPath string `toml:"-"`
}

// expandVars expands several variables with environment variables.
func (p *Profile) expandVars() (err error) {
	p.Path = strings.Replace(p.Path, "<PROFILE>", p.Name, 1)
	p.expandedPath, err = expandPath(p.Path)
	if err != nil {
		return err
	}
	p.Remote = os.ExpandEnv(p.Remote)
	p.Remote = expandRemote(p.Remote)
	return err
}

func expandPath(path string) (string, error) {
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

func expandRemote(remote string) string {
	if strings.HasPrefix(remote, "git@") || strings.HasPrefix(remote, "https://") || remote == "" {
		return remote
	}
	return "https://" + remote + ".git"
}
