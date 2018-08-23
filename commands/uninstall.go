package commands

import (
	"fmt"

	"github.com/relnod/dotm/dotfiles"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uninstalling")
		traverser := dotfiles.NewTraverser(nil)
		traverser.Traverse(Source, Destination, dotfiles.NewUnlinkAction(true))
	},
}

func init() {
	uninstallCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	uninstallCmd.Flags().StringVarP(&Destination, "destination", "d", "", "Destination directory to write to")
	rootCmd.AddCommand(uninstallCmd)
}
