package config

import (
	"os"
	"os/user"
)

// Find tries to find a dotfile config location in the home directory.
func Find() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	path := usr.HomeDir + ".dotfiles.toml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil
	}

	return usr.HomeDir, nil
}
