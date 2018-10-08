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
			cmd.Println(msgUnistallFail)
			return err
		}

		for name := range c.Profiles {
			fmt.Println(name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
