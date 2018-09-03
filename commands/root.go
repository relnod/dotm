package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dotm",
	Short: "Dotm is a dotfile manager",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes the root command.
// This is the entrypoint for the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
