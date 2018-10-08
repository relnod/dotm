package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	msgListFail = "Failed to list profiles"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all profiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := loadConfig()
		if err != nil {
			return
		}

		for name := range c.Profiles {
			fmt.Println(name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
