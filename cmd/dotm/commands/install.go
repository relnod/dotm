package commands

import (
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
	Use:   "install [remote]",
	Short: "Install the dotfiles",
	Long:  `Installs the dotfiles from a remote path. If no profile was specified, the profile name will be "default"`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		c, err := loadConfig(newFS())
		if err != nil {
			c = config.NewConfig(newFS())
		}
		c.Profiles[profile] = &config.Profile{
			Remote:   args[0],
			Path:     path,
			Excludes: excludes,
			Includes: includes,
		}

		err = dotfiles.Install(c, []string{profile}, configPath, &dotfiles.InstallOptions{
			Dry:   dry,
			Force: force,
		})
		if err != nil {
			printl(msgInstallFail, args[0])
			return err
		}

		printl(msgInstallSuccess)
		return nil
	},
}

func init() {
	installCmd.Flags().StringVar(&path, "path", "$HOME/.dotfiles/", "local git path")
	installCmd.Flags().StringVarP(&profile, "profile", "p", "default", "profile name")
	installCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwriting files")
	installCmd.Flags().BoolVar(&dry, "dry", false, "perform a dry run")
	installCmd.Flags().StringVarP(&configPath, "config", "c", "$HOME/.dotfiles.toml", "Config location")
	installCmd.Flags().StringSliceVar(&excludes, "excludes", nil, "directories to be excluded")
	installCmd.Flags().StringSliceVar(&includes, "includes", nil, "directories to be included")
	rootCmd.AddCommand(installCmd)
}
