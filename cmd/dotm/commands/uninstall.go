package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/dotfiles"
)

var (
	msgUnistallSuccess = "Dotfiles where uninstalled successfully"
	msgUnistallFail    = "Failed to uninstall dotfiles"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [profiles]",
	Short: "Uninstall the profiles",
	Long:  `Removes all symlinks for the given profiles. When profile "all" is set, all profiles will get uninstalled.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig(newFS())
		if err != nil {
			fmt.Printf("Failed to read config\n")
			return err
		}

		err = dotfiles.Uninstall(c, args, &dotfiles.UninstallOptions{
			Dry: dry,
		})
		if err != nil {
			printl(msgUnistallFail)
			return err
		}

		printl(msgUnistallSuccess)
		return nil
	},
}

func init() {
	uninstallCmd.Flags().StringVarP(&configPath, "config", "c", "$HOME/.dotfiles.toml", "config location")
	uninstallCmd.Flags().StringSliceVar(&excludes, "excludes", nil, "directories to be excluded")
	uninstallCmd.Flags().StringSliceVar(&includes, "includes", nil, "directories to be included")
	uninstallCmd.Flags().BoolVar(&dry, "dry", false, "perform a dry run")
	rootCmd.AddCommand(uninstallCmd)
}
