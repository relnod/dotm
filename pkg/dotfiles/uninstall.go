package dotfiles

import (
	"os/user"

	"github.com/relnod/dotm/pkg/config"
)

// Uninstall uninstalles the dotfiles.
func Uninstall(c *config.Config) error {
	var err error

	err = c.Validate()
	if err != nil {
		return err
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	err = UnLink(c.Path, usr.HomeDir, nil)
	if err != nil {
		return err
	}

	// err = os.RemoveAll(c.Path)
	// if err != nil {
	// 	return err
	// }

	return nil
}
