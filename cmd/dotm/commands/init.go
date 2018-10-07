package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/dotfiles"
)

var (
	remote         = ""
	msgInitSuccess = "Dotfiles where initialized successfully"
	msgInitFail    = "Failed to initialize dotfiles at '%s'"
)

var initCmd = &cobra.Command{
	Use:   "init path",
	Short: "Initialize the dotfiles",
	Long:  `Initializes the dotfiles from the given path. If no profile was specified, the profile name will be "default"`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := loadOrCreateConfig()
		c.Profiles[profile] = &config.Profile{
			Remote:   remote,
			Path:     args[0],
			Excludes: excludes,
			Includes: includes,
		}
		err := c.Profiles[profile].Initialize()
		if err != nil {
			cmd.Println(fmt.Sprintf(msgInitFail, args[0]))
			return err
		}

		err = dotfiles.Init(c, []string{profile}, os.ExpandEnv(configPath), &dotfiles.InitOptions{
			Dry:   dry,
			Force: force,
		})
		if err != nil {
			cmd.Println(fmt.Sprintf(msgInitFail, args[0]))
			return err
		}

		cmd.Println(msgInitSuccess)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&remote, "remote", "r", "", "remote git location")

	addForceFlag(initCmd)
	addBaseFlags(initCmd)

	rootCmd.AddCommand(initCmd)
}
