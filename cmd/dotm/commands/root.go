package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
)

var (
	configPath string
)

func loadConfig() (*config.Config, error) {
	var err error
	if configPath == "" {
		configPath, err = config.Find()
		if err != nil {
			return nil, err
		}
	} else {
		configPath = os.ExpandEnv(configPath)
	}
	c, err := config.NewFromTomlFile(configPath)
	if err != nil {
		return nil, err
	}

	return c, nil
}

var rootCmd = &cobra.Command{
	Use:   "dotm",
	Short: "Dotm is a dotfile manager",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "$HOME/.dotfiles.toml", "config location")
}

// Execute executes the root command.
// This is the entrypoint for the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
