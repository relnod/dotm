package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/dotfiles"
)

var (
	msgUnistallSuccess = "Dotfiles where uninstalled successfully"
	msgUnistallFail    = "Failed to uninstall dotfiles"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [profiles]",
	Short: "Uninstall the profiles",
	Long:  `Removes all symlinks for the given profiles. When profile "all" is set, all profiles will get uninstalled.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			cmd.Println(msgUnistallFail)
			return err
		}

		if len(args) == 0 {
			args = []string{"all"}
		}

		err = dotfiles.Uninstall(c, args, &dotfiles.UninstallOptions{
			Dry: dry,
		})
		if err != nil {
			cmd.Println(msgUnistallFail)
			return err
		}

		cmd.Println(msgUnistallSuccess)
		return nil
	},
	ValidArgs: []string{"$(dotm list)"},
}

func init() {
	addBaseFlags(uninstallCmd)

	rootCmd.AddCommand(uninstallCmd)
}
