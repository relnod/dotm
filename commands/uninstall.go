package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/dotfiles"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.NewFromTomlFile(configPath)
		if err != nil {
			fmt.Printf("Failed to read config\n")
			fmt.Printf("Error: %s\n", err)
			return
		}
		err = dotfiles.Uninstall(c)
		if err != nil {
			fmt.Printf("Failed to uninstall dotfiles\n")
			fmt.Printf("Error: %s\n", err.Error())
			return
		}

		fmt.Println("Dotfiles where unistalled successfully")
	},
}

func init() {
	uninstallCmd.Flags().StringVarP(&configPath, "config", "c", "", "config location")
	rootCmd.AddCommand(uninstallCmd)
}
