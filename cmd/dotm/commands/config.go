package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	configCmdRemote      bool
	configCmdPath        bool
	configCmdIncludes    bool
	configCmdExcludes    bool
	configCmdPreUpdates  bool
	configCmdPostUpdates bool
)

var configCmd = &cobra.Command{
	Use:       "config Profile",
	Short:     "Prints information about a profile",
	Long:      ``,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"$(dotm list)"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			return err
		}
		p, err := c.FindProfile(args[0])
		if err != nil {
			return err
		}

		printString := func(s string) {
			fmt.Println(s)
		}
		printStringSlice := func(r []string) {
			for _, s := range r {
				printString(s)
			}
		}

		if configCmdRemote {
			printString(p.Remote)
		}
		if configCmdPath {
			printString(p.Path)
		}
		if configCmdIncludes {
			printStringSlice(p.Includes)
		}
		if configCmdExcludes {
			printStringSlice(p.Excludes)
		}
		if configCmdPreUpdates {
			printStringSlice(p.PreUpdate)
		}
		if configCmdPostUpdates {
			printStringSlice(p.PostUpdate)
		}

		return nil
	},
}

func init() {
	configCmd.Flags().BoolVarP(&configCmdRemote, "remote", "", false, "show remote")
	configCmd.Flags().BoolVarP(&configCmdPath, "path", "", false, "show path")
	configCmd.Flags().BoolVarP(&configCmdIncludes, "includes", "", false, "show includes")
	configCmd.Flags().BoolVarP(&configCmdExcludes, "excludes", "", false, "show excludes")
	configCmd.Flags().BoolVarP(&configCmdPreUpdates, "pre-updates", "", false, "show pre update hooks")
	configCmd.Flags().BoolVarP(&configCmdPostUpdates, "post-updates", "", false, "show post update hooks")
	rootCmd.AddCommand(configCmd)
}
