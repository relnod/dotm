package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Access config of a profile",
	Long:  ``,
}

var configGetCmd = &cobra.Command{
	Use:       "get profile",
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

		if cfgGetRemote {
			printString(p.Remote)
		}
		if cfgGetPath {
			printString(p.Path)
		}
		if cfgGetHooksEnabled {
			printString(strconv.FormatBool(p.HooksEnabled))
		}
		if cfgGetIncludes {
			printStringSlice(p.Includes)
		}
		if cfgGetExcludes {
			printStringSlice(p.Excludes)
		}
		if cfgGetPreUpdates {
			printStringSlice(p.PreUpdate)
		}
		if cfgGetPostUpdates {
			printStringSlice(p.PostUpdate)
		}

		return nil
	},
}

var (
	cfgGetRemote       bool
	cfgGetPath         bool
	cfgGetHooksEnabled bool
	cfgGetIncludes     bool
	cfgGetExcludes     bool
	cfgGetPreUpdates   bool
	cfgGetPostUpdates  bool
)

var configSetCmd = &cobra.Command{
	Use:       "set profile",
	Short:     "Sets information for a profile",
	Long:      ``,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"$(dotm list)"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := dotm.LoadConfig()
		if err != nil {
			return err
		}
		p, err := c.Profile(args[0])
		if err != nil {
			return err
		}

		if cmd.Flag("remote").Changed {
			p.Remote = cfgSetRemote
		}
		if cmd.Flag("path").Changed {
			p.Path = cfgSetPath
		}
		if cmd.Flag("hooks-enabled").Changed {
			p.HooksEnabled = cfgSetHooksEnabled
		}
		if cmd.Flag("includes").Changed {
			p.Includes = cfgSetIncludes
		}
		if cmd.Flag("excludes").Changed {
			p.Excludes = cfgSetExcludes
		}
		if cmd.Flag("pre-updates").Changed {
			p.PreUpdate = cfgSetPreUpdates
		}
		if cmd.Flag("post-updates").Changed {
			p.PostUpdate = cfgSetPostUpdates
		}

		return c.Write()
	},
}

var (
	cfgSetRemote       string
	cfgSetPath         string
	cfgSetHooksEnabled bool
	cfgSetIncludes     []string
	cfgSetExcludes     []string
	cfgSetPreUpdates   []string
	cfgSetPostUpdates  []string
)

func init() {
	configGetCmd.Flags().BoolVar(&cfgGetRemote, "remote", false, "Prints the remote")
	configGetCmd.Flags().BoolVar(&cfgGetPath, "path", false, "Prints the path")
	configGetCmd.Flags().BoolVar(&cfgGetHooksEnabled, "hooks-enabled", false, "Prints, wether hooks are enabled")
	configGetCmd.Flags().BoolVar(&cfgGetIncludes, "includes", false, "Prints includes")
	configGetCmd.Flags().BoolVar(&cfgGetExcludes, "excludes", false, "Prints excludes")
	configGetCmd.Flags().BoolVar(&cfgGetPreUpdates, "pre-updates", false, "Prints pre update hooks")
	configGetCmd.Flags().BoolVar(&cfgGetPostUpdates, "post-updates", false, "Prints post update hooks")

	configSetCmd.Flags().StringVar(&cfgSetRemote, "remote", "", "Sets the remote")
	configSetCmd.Flags().StringVar(&cfgSetPath, "path", "", "Sets the path")
	configSetCmd.Flags().BoolVar(&cfgSetHooksEnabled, "hooks-enabled", false, "Enables/disables hooks")
	configSetCmd.Flags().StringSliceVar(&cfgSetIncludes, "includes", nil, "Sets includes")
	configSetCmd.Flags().StringSliceVar(&cfgSetExcludes, "excludes", nil, "Sets excludes")
	configSetCmd.Flags().StringSliceVar(&cfgSetPreUpdates, "pre-updates", nil, "Sets pre update hooks")
	configSetCmd.Flags().StringSliceVar(&cfgSetPostUpdates, "post-updates", nil, "Sets post update hooks")

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
