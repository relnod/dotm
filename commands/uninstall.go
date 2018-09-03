package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/dotfiles"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the dotfiles",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to read config\n")
			return err
		}
		err = dotfiles.Uninstall(c)
		if err != nil {
			fmt.Printf("Failed to uninstall dotfiles\n")
			return err
		}

		fmt.Println("Dotfiles where unistalled successfully")
		return nil
	},
}

func init() {
	uninstallCmd.Flags().StringVarP(&configPath, "config", "c", "", "config location")
	rootCmd.AddCommand(uninstallCmd)
}
