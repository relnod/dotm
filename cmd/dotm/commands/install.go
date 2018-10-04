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
	Use:   "install",
	Short: "Install the dotfiles",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		c := config.New(&config.Config{
			Remote:   args[0],
			Path:     path,
			FS:       newFS(),
			Excludes: *excludes,
			Includes: *includes,
		})

		err = dotfiles.Install(c, configPath)
		if err != nil {
			printl(msgInstallFail, args[0])
			return err
		}

		printl(msgInstallSuccess)
		return nil
	},
}

func init() {
	installCmd.Flags().StringVarP(&path, "path", "p", "$HOME/.dotfiles/", "Local git path")
	rootCmd.AddCommand(installCmd)
}
