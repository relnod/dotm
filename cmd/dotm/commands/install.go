package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const installHelp = `Installs dotfiles from a remote git repository.
If no profile was specified, the profile name will be "default".`

const installExamples = `dotm install github.com/user/dotfiles
dotm install --profile=myprofile github.com/user/dotfiles`

var installCmd = &cobra.Command{
	Use:     "install remote",
	Short:   "Install dotfiles from a remote git repository",
	Long:    installHelp,
	Example: installExamples,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.InstallWithContext(
			interruptContext(),
			&dotm.Profile{
				Name:   profile,
				Remote: args[0],
				Path:   sanitizePath(path, profile),
				// When installing a dotfile profile, hooks are disabled by
				// default.
				HooksEnabled: false,
			},
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
