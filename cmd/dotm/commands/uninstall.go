package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const uninstallHelp = `Removes all symlinks for the given profile.
Tries to restore backup files.

Example:
dotm uninstall default
`

var uninstallCmd = &cobra.Command{
	Use:       "uninstall profile",
	Short:     "Uninstall the profile",
	Long:      uninstallHelp,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"$(dotm list)"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.Uninstall(args[0], &dotm.UninstallOptions{
			Dry:   dry,
			Clean: clean,
		})
	},
}

var clean bool

func init() {
	uninstallCmd.Flags().BoolVar(&dry, "dry", false, "perfomes a dry run")
	uninstallCmd.Flags().BoolVar(&clean, "clean", false, "delete the local path and remove profile from config")

	rootCmd.AddCommand(uninstallCmd)
}
