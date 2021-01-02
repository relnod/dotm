package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/relnod/dotm/cmd/dotm/commands"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"dotm":   execdotm,
		"islink": execislink,
	}))
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(e *testscript.Env) error {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			path := os.Getenv("PATH")

			e.Vars = []string{
				// Set $HOME to the temporary home.
				"HOME=" + e.WorkDir + "/home/testuser",
				// $PATH is needed from host for the git and touch binary.
				"PATH=" + path,
				// OLDWD is needed to execute the dotm binary.
				"OLDWD=" + wd,
			}

			return nil
		},
	})
}

func execdotm() int {
	err := commands.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}

// execislink checks if os.Args[1] is a symlink.
//
// If os.Args[2] is given it checks if the symlink resolves to os.Args[2].
func execislink() int {
	_, err := os.Lstat(os.Args[1])
	if err != nil {
		return 1
	}
	link, err := os.Readlink(os.Args[1])
	if err != nil {
		return 1
	}
	if len(os.Args) > 2 && link != os.Args[2] {
		return 1
	}
	return 0
}
