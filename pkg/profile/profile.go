package profile

import (
	"os"
	"path/filepath"
	"strings"
)

// Profile defines one set of dotfiles.
type Profile struct {
	Remote     string   `toml:"remote"`
	Path       string   `toml:"path"`
	Includes   []string `toml:"includes"`
	Excludes   []string `toml:"excludes"`
	PreUpdate  []string `toml:"pre_update"`
	PostUpdate []string `toml:"post_update"`
}

// Initialize sets the profile up
func (c *Profile) Initialize(name string) (err error) {
	c.Remote = os.ExpandEnv(c.Remote)
	c.Path, err = expandPath(c.Path, name)

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
