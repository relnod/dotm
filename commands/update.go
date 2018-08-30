package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/dotfiles"
)

var (
	configPath string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.NewFromTomlFile(configPath)
		if err != nil {
			fmt.Printf("Failed to read config\n")
			fmt.Printf("Error: %s\n", err)
			return
		}
		err = dotfiles.Update(c)
		if err != nil {
			fmt.Printf("Failed to upate dotfiles\n")
			fmt.Printf("Error: %s\n", err.Error())
			return
		}

		fmt.Println("Dotfiles where updated successfully")
	},
}

func init() {
	updateCmd.Flags().StringVarP(&configPath, "config", "c", "", "config location")
	rootCmd.AddCommand(updateCmd)
}
