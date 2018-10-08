package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/dotfiles"
	"github.com/relnod/dotm/pkg/profile"
)

var (
	path              string
	msgInstallSuccess = "Dotfiles where installed successfully"
	msgInstallFail    = "Failed to install dotfiles from '%s'"
)

var installCmd = &cobra.Command{
	Use:   "install remote",
	Short: "Install the dotfiles",
	Long:  `Installs the dotfiles from a remote path. If no profile was specified, the profile name will be "default"`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := loadOrCreateConfig()
		err := c.AddProfile(profileName, &profile.Profile{
			Remote:   args[0],
			Path:     path,
			Excludes: excludes,
			Includes: includes,
		})
		if err != nil {
			cmd.Println(fmt.Sprintf(msgInstallFail, args[0]))
			return err
		}

		err = dotfiles.Install(c, []string{profileName}, configPath, &dotfiles.InstallOptions{
			Dry:   dry,
			Force: force,
		})
		if err != nil {
			cmd.Println(fmt.Sprintf(msgInstallFail, args[0]))
			return err
		}

		cmd.Println(msgInstallSuccess)
		return nil
	},
}

func init() {
	addPathFlag(installCmd)
	addForceFlag(installCmd)
	addBaseFlags(installCmd)

	rootCmd.AddCommand(installCmd)
}
