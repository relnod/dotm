package config

import (
	"os"
	"os/user"

	"github.com/relnod/fsa"
)

// Find tries to find a dotfile config location in the home directory.
func Find(fs fsa.FileSystem) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	path := usr.HomeDir + ".dotfiles/dotm.toml"
	if _, err := fs.Stat(path); os.IsNotExist(err) {
		return "", nil
	}

	return usr.HomeDir, nil
}
