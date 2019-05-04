package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const initHelp = `Initializes a new dotfile profile from the given path.
If no profile was set, the profile name will be "default"

Example:
dotm init --profile=myprofile $HOME/dotfiles`

var initCmd = &cobra.Command{
	Use:   "init path",
	Short: "Initialize a new dotfile profile from the given path.",
	Long:  initHelp,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.Init(
			&dotm.Profile{Name: profile, Path: args[0]},
			&dotm.InitOptions{
				LinkOptions: linkOptionsFromFlags(),
			},
		)
	},
}

func init() {
	addProfileFlag(initCmd)
	addLinkFlags(initCmd)
	rootCmd.AddCommand(initCmd)
}
