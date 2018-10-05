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
	Use:   "init [path]",
	Short: "Initialize the dotfiles",
	Long:  `Initializes the dotfiles from the given path. If no profile was specified, the profile name will be "default"`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		c, err := loadConfig(newFS())
		if err != nil {
			c = config.NewConfig(newFS())
		}
		c.Profiles[profile] = &config.Profile{
			Remote:   remote,
			Path:     path,
			Excludes: excludes,
			Includes: includes,
		}

		err = dotfiles.Init(config.New(c), []string{profile}, os.ExpandEnv(configPath), &dotfiles.InitOptions{
			Dry:   dry,
			Force: force,
		})
		if err != nil {
			printl(msgInitFail, path)
			return err
		}

		printl(msgInitSuccess)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&remote, "remote", "r", "", "remote git location")
	initCmd.Flags().StringVarP(&profile, "profile", "p", "Default", "Profile name")
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwriting files")
	initCmd.Flags().BoolVar(&dry, "dry", false, "perform a dry run")
	initCmd.Flags().StringVarP(&configPath, "config", "c", "$HOME/.dotfiles.toml", "config location")
	initCmd.Flags().StringSliceVar(&excludes, "excludes", nil, "directories to be excluded")
	initCmd.Flags().StringSliceVar(&includes, "includes", nil, "directories to be included")
	rootCmd.AddCommand(initCmd)
}
