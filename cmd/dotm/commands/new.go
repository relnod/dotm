package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/dotfiles"
	"github.com/relnod/dotm/pkg/profile"
)

var (
	msgNewSuccess = "New Profile created"
	msgNewFail    = "Failed to create new profile"
)

var newCmd = &cobra.Command{
	Use:   "new profile",
	Short: "Creates a new dotfile profile",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := loadOrCreateConfig()
		err := c.AddProfile(args[0], &profile.Profile{
			Remote:   remote,
			Path:     path,
			Excludes: excludes,
			Includes: includes,
		})
		if err != nil {
			cmd.Println(msgNewFail)
			return err
		}

		err = dotfiles.New(c, []string{args[0]}, configPath)
		if err != nil {
			cmd.Println(msgNewFail)
			return err
		}

		cmd.Println(msgNewSuccess)
		return nil
	},
}

func init() {
	addPathFlag(newCmd)
	addConfigFlag(newCmd)

	rootCmd.AddCommand(newCmd)
}
