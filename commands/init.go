package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/dotfiles"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize dotfiles",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		c := &config.Config{
			Remote: remote,
			Path:   destination,
		}

		err = dotfiles.Init(c)
		if err != nil {
			fmt.Printf("Failed to install dotfiles from '%s'\n", remote)
			return err
		}

		fmt.Println("Dotfiles where installed successfully")
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&configPath, "config", "c", "", "config location")
	initCmd.Flags().StringVarP(&destination, "destination", "d", "~/.dotfiles2/", "Local git destination")
	rootCmd.AddCommand(initCmd)
}
