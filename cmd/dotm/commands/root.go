package commands

import (
	"os"

	"github.com/relnod/fsa"
	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
)

var (
	configPath string
	profile    string
	testRoot   string
	excludes   *[]string
	includes   *[]string
)

func newFS() (fs fsa.FileSystem) {
	fs = fsa.NewOsFs()
	if testRoot != "" {
		fs = fsa.NewBaseFs(fs, testRoot)
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
	Long: `Dotm is a dotfile manager. It works by symlinking the files from the dotfile folder to its corresponding place under the home directory of the user.

Configuration file
The configuration file is located at $HOME/.dotfiles.toml (can be changed with the --config flag). It can hold multiple profiles. Each profile consists of a path to the local dotfile location and an optional remote path to a git repository.

Example:
toml
# You can define multiple profiles
[profiles.default]

# Upstream git repository
remote = "github.com/relnod/dotm"

# Path to local git repository
path = ".dotfiles/"

# Configs to be included
includes = [
    "bash",
    "nvim",
    "tmux"
]

# Configs to be excluded
excludes = [
    "bash",
    "nvim",
    "tmux"
]

Dotfiles folder
A Dotfile folder consists of multiple top level directories to group similar configuration files (e.g. "vim" or "tmux"). The file structure below those top level directories are directly mapped to $HOME.

Example:
tmux/.tmux.conf             -> $HOME/.tmux.conf
bash/.bashrc                -> $HOME/.bashrc
nvim/.config/nvim/init.vim  -> $HOME/.config/nvim/init.vim

From an existing dotfile folder:
dotm init <path-to-existing-dotfile-folder>

From a remote git repository:
dotm install <url-to-remote-repository>
	`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "$HOME/.dotfiles.toml", "config location")
	excludes = rootCmd.PersistentFlags().StringSlice("excludes", nil, "Directories to be excluded")
	includes = rootCmd.PersistentFlags().StringSlice("includes", nil, "Directories to be included")
	rootCmd.PersistentFlags().StringVarP(&testRoot, "testRoot", "", "", "root location (used for testing puposes)")
	rootCmd.PersistentFlags().MarkHidden("testRoot")
}

// Execute executes the root command.
// This is the entrypoint for the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
