package commands

import (
	"os"

	"github.com/relnod/fsa"
	"github.com/spf13/cobra"

	"github.com/relnod/dotm/pkg/config"
)

var (
	genCompletions bool
	testRoot       string
)

func newFS() (fs fsa.FileSystem) {
	fs = fsa.NewOsFs()
	if testRoot != "" {
		fs = fsa.NewBaseFs(fs, testRoot)
	}
	return fs
}

func loadConfig() (*config.Config, error) {
	var err error
	fs := newFS()
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

func loadOrCreateConfig() *config.Config {
	c, err := loadConfig()
	if err != nil {
		c = config.NewConfig(newFS())
	}
	return c
}

var rootCmd = &cobra.Command{
	Use:   "dotm",
	Short: "Dotm is a dotfile manager",
	Long: `Dotm is a dotfile manager. It works by symlinking the files from the dotfile folder to its corresponding place under the home directory of the user.

Configuration file
The configuration file is located at $HOME/.dotfiles/dotm.toml (can be changed with the --config flag). It can hold multiple profiles. Each profile consists of a path to the local dotfile location and an optional remote path to a git repository.

Example:
# You can define multiple profiles
[profiles.default]

# Upstream git repository
remote = "github.com/relnod/dotm"

# Path to local git repository
path = ".dotfiles/default/"

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

pre_update = [
    "echo 'pre update'"
]

post_update = [
    "echo 'post update'"
]

Dotfiles folder
A dotfile folder consists of multiple top level directories to group similar configuration files (e.g. "vim" or "tmux"). Includes and excludes of those top level directories can be defined in the config file or with flags (--exclude, --includes). Also all top level directories with the prefix "_" will be excluded. The file structure below those top level directories is directly mapped to $HOME.

Example:
tmux/.tmux.conf             -> $HOME/.tmux.conf
bash/.bashrc                -> $HOME/.bashrc
nvim/.config/nvim/init.vim  -> $HOME/.config/nvim/init.vim

From an existing dotfile folder:
dotm init <path-to-existing-dotfile-folder>

From a remote git repository:
dotm install <url-to-remote-repository>

Hooks
Update hooks can be applied via global config, at profile root and per top level directory. For hooks at profile root and top level directory you can create a hooks.toml. Note: This file won't be symlinked.

Example:
pre_update = [
    "nvim +PlugInstall +qall"
]
	`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if genCompletions {
			return cmd.GenBashCompletion(os.Stdout)
		}
		return cmd.Usage()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&testRoot, "testRoot", "", "", "root location (used for testing purposes)")
	rootCmd.PersistentFlags().MarkHidden("testRoot")

	rootCmd.Flags().BoolVarP(&genCompletions, "genCompletions", "", false, "generate bash completions")
}

// Execute executes the root command.
// This is the entrypoint for the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
