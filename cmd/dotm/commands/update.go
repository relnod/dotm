package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/dotfiles"
)

var (
	updateFromRemote bool
	msgUpdateSuccess = "Dotfiles were updated successfully"
	msgUpdateFail    = "Failed to update dotfiles"
)

var updateCmd = &cobra.Command{
	Use:   "update [profiles]",
	Short: "Update the profiles",
	Long:  `Updates all symlinks for the given profiles. When profile "all" is set, all profiles will get updated.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig(newFS())
		if err != nil {
			fmt.Printf("Failed to read config\n")
			return err
		}
		err = dotfiles.Update(c, args, &dotfiles.UpdateOptions{
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
	updateCmd.Flags().BoolVar(&updateFromRemote, "fromRemote", false, "update from remote")
	updateCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwriting files")
	updateCmd.Flags().StringVarP(&configPath, "config", "c", "$HOME/.dotfiles.toml", "config location")
	updateCmd.Flags().StringSliceVar(&excludes, "excludes", nil, "directories to be excluded")
	updateCmd.Flags().StringSliceVar(&includes, "includes", nil, "directories to be included")
	rootCmd.AddCommand(updateCmd)
}
