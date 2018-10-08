package commands

import (
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
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			cmd.Println(msgUpdateFail)
			return err
		}

		if len(args) == 0 {
			args = []string{"all"}
		}

		err = dotfiles.Update(c, args, &dotfiles.UpdateOptions{
			UpdateFromRemote: updateFromRemote,
			Force:            force,
			Dry:              dry,
		})
		if err != nil {
			cmd.Println(msgUpdateFail)
			return err
		}

		cmd.Println(msgUpdateSuccess)
		return err
	},
	ValidArgs: []string{"$(dotm list)"},
}

func init() {
	updateCmd.Flags().BoolVar(&updateFromRemote, "fromRemote", false, "update from remote")

	addForceFlag(updateCmd)
	addBaseFlags(updateCmd)

	rootCmd.AddCommand(updateCmd)
}
