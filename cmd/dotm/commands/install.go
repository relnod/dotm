package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const installHelp = `Installs dotfiles from a remote git repository.
If no profile was specified, the profile name will be "default".

Example:
dotm install --profile=myprofile github.com/user/dotfiles`

var installCmd = &cobra.Command{
	Use:   "install remote",
	Short: "Install dotfiles from a remote git repository",
	Long:  installHelp,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.Install(
			&dotm.Profile{Name: profile, Remote: args[0], Path: path},
			&dotm.InstallOptions{
				LinkOptions: linkOptionsFromFlags(),
			},
		)
	},
}

func init() {
	addPathFlag(installCmd)
	addProfileFlag(installCmd)
	addLinkFlags(installCmd)
	rootCmd.AddCommand(installCmd)
}
