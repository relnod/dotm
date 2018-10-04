package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
)

// Uninstall uninstalles the dotfiles.
func Uninstall(c *config.Config, names []string) error {
	profiles, err := c.FindProfiles(names...)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		err = UnLinkProfile(c.FS, p)
		if err != nil {
			return err
		}
		// err = os.RemoveAll(c.Path)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}
