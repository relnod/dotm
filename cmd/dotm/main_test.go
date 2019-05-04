package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/relnod/dotm/cmd/dotm/commands"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"dotm":   execdotm,
		"islink": execislink,
		"debug":  execdebug,
	}))
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(e *testscript.Env) error {
			// PATH is needed from host for the git binary.
			e.Vars = []string{"HOME=" + e.WorkDir + "/home/testuser", "PATH=/usr/bin/"}
			return nil
		},
	})
}

func execdotm() int {
	err := commands.Execute()
	if err != nil {
		return 1
	}
	return 0
}

func execislink() int {
	_, err := os.Lstat(os.Args[1])
	if err != nil {
		return 1
	}
	_, err = os.Readlink(os.Args[1])
	if err != nil {
		return 1
	}
	return 0
}

func execdebug() int {
	cmd := exec.Command("ls", "-la", os.ExpandEnv("$HOME"))
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return 1
	}
	return 0
}
