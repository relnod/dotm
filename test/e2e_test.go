package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/relnod/fsa"
	"github.com/relnod/fsa/testutil"
	"github.com/stretchr/testify/assert"
)

var coverage bool

func preTest() error {
	if os.Getenv("ROOT") == "" {
		return fmt.Errorf("$ROOT should be set, but is empty")
	}

	out, err := execCommand(`
cd $ROOT/test/testdata/remote && \
    git init && \
    git config --local user.email "you@example.com" && \
    git config --local user.name "Your Name" && \
    git add . && \
    git commit -m "initial commit"
`)
	if err != nil {
		fmt.Println(out)
		return err
	}

	if coverage {
		out, err = execCommand("rm $ROOT/coverage-e2e-*")
		if err != nil {
			fmt.Println(out)
			// return err
		}
	}
	return nil
}

func postTest() error {
	out, err := execCommand("rm -rf $ROOT/test/testdata/remote/.git")
	if err != nil {
		fmt.Println(out)
		return err
	}

	if coverage {
		out, err = execCommand("gocovmerge $ROOT/coverage-e2e-* > $ROOT/coverage-e2e.out")
		if err != nil {
			fmt.Println(out)
			// return err
		}
		out, err = execCommand("rm $ROOT/coverage-e2e-*")
		if err != nil {
			fmt.Println(out)
			// return err
		}
	}
	return nil
}

func runTests(t *testing.T, cmd string) {
	if err := preTest(); err != nil {
		t.Fatal(err)
	}
	defer postTest()

	testcases, err := parseDir("./testcases")
	if err != nil {
		t.Fatal(err)
	}

	for i, c := range testcases {
		t.Run(c.name, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			c.given = strings.Replace(c.given, "$BASE", fs.Base(), -1)
			assert.NoError(tt, testutil.CreateFiles(fs, c.given))
			assert.NoError(tt, testutil.AddFiles(fs, "./testdata", "/"))
			out, err := c.exec(cmd, fs.Base(), i)
			assert.NoError(tt, err)
			if !coverage && c.cmdOutput != "" {
				assert.Equal(tt, out, c.cmdOutput)
			}
			assert.NoError(tt, testutil.CheckFiles(fs, c.expected))
		})
	}
}

func TestNormal(t *testing.T) {
	runTests(t, os.ExpandEnv("$ROOT")+"build/dotm")
}

func TestCoverage(t *testing.T) {
	coverage = true
	runTests(t, os.ExpandEnv("$ROOT")+"build/dotm.test -test.run='^TestRunMain'")
}

func TestDocker(t *testing.T) {
	if os.Getenv("SKIP_TEST_DOCKER") != "" {
		t.Skipf("Reason: %s", os.Getenv("SKIP_TEST_DOCKER"))
	}
	runTests(t, "docker run -v /tmp:/tmp --env USER=${USER} relnod/dotm:latest")
}
