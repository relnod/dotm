// This file includes a set of reusable flags.

package commands

import "github.com/spf13/cobra"

var (
	configPath string
	profile    string
	force      bool
	dry        bool
	excludes   []string
	includes   []string
)

func addConfigFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&configPath, "config", "c", "$HOME/.dotfiles/dotm.toml", "config location")
}

func addIncludeExcludeFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceVar(&excludes, "excludes", nil, "directories to be excluded")
	cmd.Flags().StringSliceVar(&includes, "includes", nil, "directories to be included")
}

func addForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&force, "force", "f", false, "force overwriting files")
}

func addDryFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&dry, "dry", false, "perform a dry run")
}

func addProfileFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&profile, "profile", "p", "default", "Profile name")
}

func addBaseFlags(cmd *cobra.Command) {
	addConfigFlag(cmd)
	addProfileFlag(cmd)

	addDryFlag(cmd)
	addIncludeExcludeFlags(cmd)
}
