package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
)

// Uninstall uninstalles the dotfiles.
func Uninstall(c *config.Config) error {
	var err error

	err = c.Validate()
	if err != nil {
		return err
	}

	err = UnLink(c.Path, "/tmp/bla2", nil)
	if err != nil {
		return err
	}

	// err = os.RemoveAll(c.Path)
	// if err != nil {
	// 	return err
	// }

	return nil
}
