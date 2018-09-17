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
	Use:   "uninstall",
	Short: "Uninstall the dotfiles",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to read config\n")
			return err
		}
		err = dotfiles.Uninstall(c)
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
