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
	Use:       "uninstall [profiles]",
	Short:     "Uninstall the profiles",
	Long:      uninstallHelp,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"$(dotm list)"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := dotm.Uninstall(args[0], &dotm.UninstallOptions{
			Dry: dry,
		})
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	uninstallCmd.Flags().BoolVar(&dry, "dry", false, "perfomes a dry run")

	rootCmd.AddCommand(uninstallCmd)
}
