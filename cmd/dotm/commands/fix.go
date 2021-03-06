package commands

import (
	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

const fixHelp = `Tries to fix the configuration file. This command can be used
after upgrading dotm to fix potential breaking changes in the configuration
file.

This should the reduce the friction when upgrading dotm.

List of things that get fixed:
  - [0.3.0] move config from olf location at $HOME/.dotfiles/dotm.toml
  - [0.3.0] set ignore_prefix to "_", when not set
  - [0.4.0] set hooks_enabled to true, when not set`

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Tries to fix the configuration file",
	Long:  fixHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		return dotm.Fix()
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
