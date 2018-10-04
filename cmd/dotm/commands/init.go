package commands

import (
	"os"
	"path/filepath"

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
	Use:   "init",
	Short: "Initialize the dotfiles",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		c := config.New(&config.Config{
			Remote:   remote,
			Path:     path,
			FS:       newFS(),
			Excludes: *excludes,
			Includes: *includes,
		})

		err = dotfiles.Init(c, os.ExpandEnv(configPath))
		if err != nil {
			printl(msgInitFail, c.Path)
			return err
		}

		printl(msgInitSuccess)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&remote, "remote", "r", "", "remote git location")
	rootCmd.AddCommand(initCmd)
}
