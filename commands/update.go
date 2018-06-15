package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("updating")
		// TODO: implement
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
