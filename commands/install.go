package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/dotfiles"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("installing")
		traverser := dotfiles.NewTraverser(nil)
		traverser.Traverse(Source, Destination, dotfiles.NewLinkAction(true))
	},
}

var Source string
var Destination string

func init() {
	installCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	installCmd.Flags().StringVarP(&Destination, "destination", "d", "~", "Destination directory to write to (used for debugging)")
	rootCmd.AddCommand(installCmd)
}
