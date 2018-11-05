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
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			fmt.Println(msgInitFail)
			return err
		}

		for name := range c.Profiles {
			fmt.Println(name)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
