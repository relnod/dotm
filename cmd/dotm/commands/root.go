package commands

import (
	"os"

	"github.com/relnod/fsa"
	"github.com/relnod/fsa/basefs"
	"github.com/relnod/fsa/osfs"
	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
)

var (
	configPath string
	testRoot   string
)

func newFS() (fs fsa.FileSystem) {
	fs = osfs.New()
	if testRoot != "" {
		fs = basefs.New(fs, testRoot)
	}
	return fs
}

func loadConfig(fs fsa.FileSystem) (*config.Config, error) {
	var err error
	if configPath == "" {
		configPath, err = config.Find(fs)
		if err != nil {
			return nil, err
		}
	} else {
		configPath = os.ExpandEnv(configPath)
	}
	c, err := config.NewFromFile(fs, configPath)
	if err != nil {
		return nil, err
	}

	c.FS = fs
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
	rootCmd.PersistentFlags().StringVarP(&testRoot, "testRoot", "", "", "root location (used for testing puposes)")
}

// Execute executes the root command.
// This is the entrypoint for the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
