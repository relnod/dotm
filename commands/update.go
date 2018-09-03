package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/dotfiles"
)

var (
	configPath       string
	updateFromRemote bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the dotfiles",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to read config\n")
			return err
		}
		err = dotfiles.Update(c, &dotfiles.UpdateOptions{
			UpdateFromRemote: updateFromRemote,
		})
		if err != nil {
			fmt.Printf("Failed to upate dotfiles\n")
			return err
		}

		fmt.Println("Dotfiles where updated successfully")
		return err
	},
}

func init() {
	updateCmd.Flags().StringVarP(&configPath, "config", "c", "", "config location")
	updateCmd.Flags().BoolVar(&updateFromRemote, "remote", false, "update from remote")
	rootCmd.AddCommand(updateCmd)
}
