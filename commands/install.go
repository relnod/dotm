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
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to read config\n")
			return err
		}

		if c == nil {
			c = &config.Config{
				Remote: remote,
				Path:   destination,
			}
		}

		err = dotfiles.Install(c)
		if err != nil {
			fmt.Printf("Failed to install dotfiles from '%s'\n", remote)
			return err
		}

		fmt.Println("Dotfiles where installed successfully")
		return nil
	},
}

func init() {
	installCmd.Flags().StringVarP(&remote, "remote", "r", "", "Remote git repository")
	installCmd.Flags().StringVarP(&destination, "destination", "d", "~/.dotfiles2/", "Local git destination")
	rootCmd.AddCommand(installCmd)
}
