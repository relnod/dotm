package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const newHelp = `Creates a new dotfile profile.
Tries to initialize a new git repository.`

const newExamples = `dotm new myprofile`

var newCmd = &cobra.Command{
	Use:     "new profile",
	Short:   "Create a new dotfile profile",
	Long:    newHelp,
	Example: newExamples,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.New(
			&dotm.Profile{
				Name:         args[0],
				Path:         sanitizePath(path, args[0]),
				HooksEnabled: true,
			},
		)
	},
}

func init() {
	addPathFlag(newCmd)
	rootCmd.AddCommand(newCmd)
}
