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

	// expandedPath contains the Path after it was expanded.
	expandedPath string `toml:"-"`
}

// sanitize changes public fields in the profile. Therfore this changes the
// profile.
func (p *Profile) sanitize() {
	p.Remote = sanitizeRemote(p.Remote)
}

// sanitizeRemote checks if the remote is a valid git remote. If not it assumes
// the remote is of the form "domain/user/repo" and converts this to a valid
// https remote.
func sanitizeRemote(remote string) string {
	if strings.HasPrefix(remote, "git@") || strings.HasPrefix(remote, "https://") || remote == "" {
		return remote
	}
	return "https://" + remote + ".git"
}

// expandEnv expands several variables with environment variables.
func (p *Profile) expandEnv() (err error) {
	p.expandedPath, err = expandPath(p.Path)
	if err != nil {
		return err
	}
	return err
}

// expandPath expands the given path with environment variables and converts it
// to an absolute path, if the path is relative.
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
