package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/dotfiles"
)

var (
	remote      string
	destination string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the dotfiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := &config.Config{
			Remote: remote,
			Path:   destination,
		}

		err := dotfiles.Install(c)
		if err != nil {
			fmt.Printf("Failed to install dotfiles from '%s'\n", remote)
			fmt.Printf("Error: '%s'\n", err.Error())
			return
		}

		fmt.Println("Dotfiles where installed successfully")
	},
}

func init() {
	installCmd.Flags().StringVarP(&remote, "remote", "r", "", "Remote git repository")
	installCmd.Flags().StringVarP(&destination, "destination", "d", "~/.dotfiles2/", "Local git destination")
	rootCmd.AddCommand(installCmd)
}
