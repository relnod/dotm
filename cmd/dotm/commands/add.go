package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const addHelp = `Adds a new or existing file to the given profile.
It expectes three arguments. The profile, top level directory and the file path.
If the given file already exists under $HOME/path , it will be moved to the
profile data dir and then linked.
If the file does not exist, a new file will be created and linked.`

const addExamples = `dotm add myprofile bash .bashrc`

var addCmd = &cobra.Command{
	Use:     "add profile dir path",
	Short:   "Add a new/existing file to the profile",
	Long:    addHelp,
	Example: addExamples,
	Args:    cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.Add(
			args[0],
			args[1],
			os.ExpandEnv(args[2]),
		)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
