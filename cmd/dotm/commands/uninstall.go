package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const uninstallHelp = `Removes all symlinks for the given profile.
Tries to restore backup files.`

const uninstallExamples = `dotm uninstall default`

var uninstallCmd = &cobra.Command{
	Use:       "uninstall profile",
	Short:     "Uninstall the profile",
	Long:      uninstallHelp,
	Example:   uninstallExamples,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"$(dotm list)"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.Uninstall(args[0], &dotm.UninstallOptions{
			Dry: dry,
		})
	},
}

func init() {
	uninstallCmd.Flags().BoolVar(&dry, "dry", false, "perfomes a dry run")

	rootCmd.AddCommand(uninstallCmd)
}
