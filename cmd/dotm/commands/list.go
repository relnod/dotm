package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := dotm.LoadConfig()
		if err != nil {
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
