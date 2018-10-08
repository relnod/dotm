package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/dotfiles"
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
		err := c.AddProfile(profile, &config.Profile{
			Remote:   args[0],
			Path:     strings.Replace(path, "<PROFILE>", profile, 1),
			Excludes: excludes,
			Includes: includes,
		})
		if err != nil {
			cmd.Println(fmt.Sprintf(msgInstallFail, args[0]))
			return err
		}

		err = dotfiles.Install(c, []string{profile}, os.ExpandEnv(configPath), &dotfiles.InstallOptions{
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
	installCmd.Flags().StringVar(&path, "path", "$HOME/.dotfiles/<PROFILE>/", "local git path")

	addForceFlag(installCmd)
	addBaseFlags(installCmd)

	rootCmd.AddCommand(installCmd)
}
