package commands

import (
	"github.com/relnod/dotm"
	"github.com/spf13/cobra"
)

var path string

func addPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&path, "path", "$HOME/.config/dotm/profiles/<PROFILE>/", "path to local dotfiles profile")
}

var profile string

func addProfileFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&profile, "profile", "default", "dotfiles profile")
}

var (
	includes []string
	excludes []string
	force    bool
	dry      bool
)

func addLinkFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&force, "force", false, "overrides destination files")
	cmd.Flags().BoolVar(&dry, "dry", false, "perfomes a dry run")
	cmd.Flags().StringSliceVar(&excludes, "excludes", nil, "directories to be excluded")
	cmd.Flags().StringSliceVar(&includes, "includes", nil, "directories to be included")
}

func linkOptionsFromFlags() dotm.LinkOptions {
	return dotm.LinkOptions{
		Dry:   dry,
		Force: force,
		TraversalOptions: dotm.TraversalOptions{
			Includes: includes,
			Excludes: excludes,
		},
	}
}
