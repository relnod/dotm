package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uninstalling")
		// TODO: implement
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
