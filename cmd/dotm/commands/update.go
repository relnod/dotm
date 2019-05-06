package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const updateHelp = `Updates the symlinks for a given profile.
When the --fromRemote flag was set it first pulls from the remote git repository.
Unless the --no-hooks flag was set, pre and post update hooks are executed.

Example:
dotm update myprofile`

var updateCmd = &cobra.Command{
	Use:       "update profile",
	Short:     "Updates the symlinks for a given profile.",
	Long:      updateHelp,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"$(dotm list)"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.UpdateWithContext(
			interruptContext(),
			args[0],
			&dotm.UpdateOptions{
				FromRemote:  fromRemote,
				ExecHooks:   !noHooks,
				LinkOptions: linkOptionsFromFlags(),
			},
		)
	},
}

var (
	fromRemote bool
	noHooks    bool
)

func init() {
	updateCmd.Flags().BoolVar(&fromRemote, "fromRemote", false, "pull updates from remote")
	updateCmd.Flags().BoolVar(&noHooks, "no-hooks", false, "doesn't exec hooks")
	addLinkFlags(updateCmd)
	rootCmd.AddCommand(updateCmd)
}
