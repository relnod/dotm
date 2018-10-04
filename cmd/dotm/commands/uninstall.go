package commands

import (
	"fmt"

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
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig(newFS())
		if err != nil {
			fmt.Printf("Failed to read config\n")
			return err
		}

		err = dotfiles.Uninstall(c, args)
		if err != nil {
			printl(msgUnistallFail)
			return err
		}

		printl(msgUnistallSuccess)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
