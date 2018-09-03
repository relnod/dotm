package commands

import (
	"fmt"
	"os"

	"github.com/relnod/dotm/pkg/config"
	"github.com/spf13/cobra"
)

func loadConfig() (*config.Config, error) {
	var err error
	if configPath == "" {
		configPath, err = config.Find()
		if err != nil {
			return nil, err
		}
	}
	c, err := config.NewFromTomlFile(configPath)
	if err != nil {
		return nil, err
	}

	return c, nil
}

var rootCmd = &cobra.Command{
	Use:   "dotm",
	Short: "Dotm is a dotfile manager",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes the root command.
// This is the entrypoint for the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
