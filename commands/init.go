package commands

import (
	"fmt"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/dotfiles"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		c := &config.Config{
			Remote: remote,
			Path:   destination,
		}

		err = dotfiles.Init(c)
		if err != nil {
			fmt.Printf("Failed to install dotfiles from '%s'\n", remote)
			fmt.Printf("Error: '%s'\n", err.Error())
			return
		}

		fmt.Println("Dotfiles where installed successfully")
	},
}

func init() {
	initCmd.Flags().StringVarP(&configPath, "config", "c", "", "config location")
	initCmd.Flags().StringVarP(&destination, "destination", "d", "~/.dotfiles2/", "Local git destination")
	rootCmd.AddCommand(initCmd)
}
