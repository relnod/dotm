package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/relnod/dotm/cmd/dotm/commands"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"dotm":          execdotm,
		"islink":        execislink,
		"dotminterrupt": execdotminterrupt,
	}))
}

func TestScripts(t *testing.T) {

	// Build the dotm binary
	cmd := exec.Command("go", "build", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		t.Fatal(err)
	}

	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(e *testscript.Env) error {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			e.Vars = []string{
				// Set $HOME to the temporary home.
				"HOME=" + e.WorkDir + "/home/testuser",
				// $PATH is needed from host for the git binary.
				"PATH=/usr/bin/",
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
		return 1
	}
	return 0
}

func execdotminterrupt() int {
	cmd := exec.Command("./dotm", os.Args[1:]...)
	cmd.Dir = os.Getenv("OLDWD")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	for {
		select {
		case <-time.After(1 * time.Second):
			cmd.Process.Signal(os.Interrupt)
		case err := <-done:
			if err != nil {
				return 1
			}
			return 0
		}
	}
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
