package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const updateHelp = `Updates the symlinks for a given profile.
When the --fromRemote flag was set it first pulls from the remote git repository.
Unless the --no-hooks flag was set, pre and post update hooks are executed.

When profile is empty, all profiles get updated.`

const updateExamples = `dotm update
dotm update --force myprofile
dotm update --fromRemote myprofile`

var updateCmd = &cobra.Command{
	Use:       "update [profile]",
	Short:     "Updates the symlinks for a given profile.",
	Long:      updateHelp,
	Example:   updateExamples,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"$(dotm list)"},
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := ""
		if len(args) > 0 {
			profile = args[0]
		}
		return dotm.UpdateWithContext(
			interruptContext(),
			profile,
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
