package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute executes the root command. Returns an error on failure.
func Execute() error {
	return rootCmd.Execute()
}

const rootHelp = `dotm is a dotfile manager. It works by symlinking the files from a
version controlled dotfile folder to its corresponding place under the $HOME
directory of the user.

Configuration file
The configuration file is located at $HOME/.config/dotm/config.toml. A dotm
configuration consists of one or multiple profiles. Each profile consists of a
file path pointing to the dotfiles associated with the profile. A remote path
can be declared pointing to a remote git repository.

Example:
# You can define multiple profiles
[profiles.default]
# Path to local git repository
path = ".dotfiles/default/"

# Remote git repository
remote = "github.com/relnod/dotm"

# Top level folders to be included.
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

# Pre update hooks
pre_update = [
    "echo 'pre update'"
]

# Post update hooks
post_update = [
    "echo 'post update'"
]

Dotfiles folder

A Dotfile folder consists of multiple top level directories to group similar
configuration files (e.g. "vim" or "tmux"). The file structure below those top
level directories are directly mapped to the $HOME directory.

Example:

tmux/.tmux.conf             -> $HOME/.tmux.conf
bash/.bashrc                -> $HOME/.bashrc
nvim/.config/nvim/init.vim  -> $HOME/.config/nvim/init.vim

Hooks
Update hooks can be applied via global config, at profile root and per top level
directory. For hooks at profile root and top level directory you can create a
hooks.toml. Note: This file won't be symlinked.

Example:

# $HOME/.config/dotm/profiles/myprofile/nvim/hooks.toml
pre_update = [
    "nvim +PlugInstall +qall"
]

Usage

# New (empty) dotfile profile
dotm new myprofile

# New profile from an existing dotfile folder
dotm init <path-to-existing-dotfile-folder>
dotm init --profile=myprofile <path-to-existing-dotfile-folder>

# New profile from a remote git repository
dotm install <url-to-remote-repository>
dotm install --profile=myprofile <url-to-remote-repository>

# Updating a profile
dotm update myprofile
dotm update myprofile --fromRemote`

const bashChangeDirectory = `
function dcd {
    cd "$(dotm config "$1" --path)" || exit
}
_dcd_completions()
{
    if [ "${#COMP_WORDS[@]}" != "2" ]; then
        return
    fi
    COMPREPLY=($(compgen -W "$(dotm list)" "${COMP_WORDS[1]}"))
}
complete -F _dcd_completions dcd
`

var rootCmd = &cobra.Command{
	Use:     "dotm",
	Short:   "Dotm is a dotfile manager",
	Long:    rootHelp,
	Version: "v0.3.0",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if genCompletions {
			return cmd.GenBashCompletion(os.Stdout)
		}
		if genChangeDirectory {
			fmt.Println(bashChangeDirectory)
			return nil
		}
		return cmd.Usage()
	},
}

var (
	genCompletions     bool
	genChangeDirectory bool
)

func init() {
	rootCmd.Flags().BoolVarP(&genCompletions, "genCompletions", "", false, "generate bash completions")
	rootCmd.Flags().BoolVarP(&genChangeDirectory, "genChangeDirectory", "", false, "generate bash change directory command")
}
