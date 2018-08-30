package e2e

import (
	"testing"

	"github.com/relnod/dotm/internal/testutil"
	"github.com/relnod/dotm/internal/testutil/assert"
	"github.com/relnod/dotm/test/runner"
)

func PTestSmoke(t *testing.T) {
	subCmds := []runner.DotmCmd{
		runner.DotmCmd{
			SubCommand: "install",
			Params:     make(map[string]string),
		},
		runner.DotmCmd{
			SubCommand: "uninstall",
			Params:     make(map[string]string),
		},
	}

	for _, subCmd := range subCmds {
		t.Run(subCmd.SubCommand, func(tt *testing.T) {
			fs := testutil.NewFileSystem()
			defer fs.Cleanup()
			fs.MkdirAll("dotiles")
			fs.MkdirAll("home")
			fs.Create("dotiles/a.txt")

			subCmd.Params["source"] = fs.BasePath()
			subCmd.Params["destination"] = fs.Path("home")

			r := runner.Run(subCmd)
			assert.ErrorIsNil(tt, r.Error())
		})
	}
}
