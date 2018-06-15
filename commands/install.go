package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("installing")
		// TODO: implement
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
