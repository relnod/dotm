package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/dotfiles"
)

var (
	updateFromRemote bool
	msgUpdateSuccess = "Dotfiles where updated successfully"
	msgUpdateFail    = "Failed to update dotfiles"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the dotfiles",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig(newFS())
		if err != nil {
			fmt.Printf("Failed to read config\n")
			return err
		}
		err = dotfiles.Update(c, &dotfiles.UpdateOptions{
			UpdateFromRemote: updateFromRemote,
		})
		if err != nil {
			printl(msgUpdateFail)
			return err
		}

		printl(msgUpdateSuccess)
		return err
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateFromRemote, "remote", false, "update from remote")
	rootCmd.AddCommand(updateCmd)
}
