package main

import (
	"os"

	"github.com/relnod/dotm/cmd/dotm/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
