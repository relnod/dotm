package commands

import (
	"fmt"

	"github.com/relnod/dotm"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:       "config profile",
	Short:     "Prints information about a profile",
	Long:      ``,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"$(dotm list)"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := dotm.LoadConfig()
		if err != nil {
			return err
		}
		p, err := c.Profile(args[0])

		printString := func(s string) {
			fmt.Println(s)
		}
		printStringSlice := func(r []string) {
			for _, s := range r {
				printString(s)
			}
		}

		if cfgRemote {
			printString(p.Remote)
		}
		if cfgPath {
			printString(p.Path)
		}
		if cfgIncludes {
			printStringSlice(p.Includes)
		}
		if cfgExcludes {
			printStringSlice(p.Excludes)
		}
		if cfgPreUpdates {
			printStringSlice(p.PreUpdate)
		}
		if cfgPostUpdates {
			printStringSlice(p.PostUpdate)
		}

		return nil
	},
}

var (
	cfgRemote      bool
	cfgPath        bool
	cfgIncludes    bool
	cfgExcludes    bool
	cfgPreUpdates  bool
	cfgPostUpdates bool
)

func init() {
	configCmd.Flags().BoolVarP(&cfgRemote, "remote", "", false, "show remote")
	configCmd.Flags().BoolVarP(&cfgPath, "path", "", false, "show path")
	configCmd.Flags().BoolVarP(&cfgIncludes, "includes", "", false, "show includes")
	configCmd.Flags().BoolVarP(&cfgExcludes, "excludes", "", false, "show excludes")
	configCmd.Flags().BoolVarP(&cfgPreUpdates, "pre-updates", "", false, "show pre update hooks")
	configCmd.Flags().BoolVarP(&cfgPostUpdates, "post-updates", "", false, "show post update hooks")
	rootCmd.AddCommand(configCmd)
}
