package main

import (
	"fmt"
	"os"

	"github.com/relnod/dotm/cmd/dotm/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
