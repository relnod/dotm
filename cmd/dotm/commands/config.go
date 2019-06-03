package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/relnod/dotm"
	"github.com/relnod/dotm/internal/clic"
)

const configExamples = `dotm config ignore_prefix
dotm config ignore_prefix "_"
dotm config profile.default.path
dotm config profile.default.path "mypath"`

var configCmd = &cobra.Command{
	Use:       "config accessor [value]",
	Short:     "Sets/Gets values from the configuration file.",
	Long:      ``,
	Example:   configExamples,
	Args:      cobra.MaximumNArgs(2),
	ValidArgs: []string{"dotm config --args"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := dotm.LoadConfig()
		if err != nil {
			return err
		}
		if showArgs {
			for _, a := range clic.Args(c) {
				fmt.Println(a)
			}
			return nil
		}
		out, err := clic.Run(strings.Join(args, " "), c)
		if err != nil {
			return err
		}
		if out != "" {
			fmt.Println(out)
			return nil
		}
		return c.Write()
	},
}

var showArgs bool

func init() {
	configCmd.Flags().BoolVar(&showArgs, "args", false, "print all possible args")
	rootCmd.AddCommand(configCmd)
}
