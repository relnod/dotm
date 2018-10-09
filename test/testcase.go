package test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type testcase struct {
	name     string
	given    string
	cmd      string
	expected string
}

func (t *testcase) exec(cmd, dir string, index int) error {
	if coverage {
		cmd += " -test.coverprofile=$ROOT/coverage-e2e-" + strconv.Itoa(index) + ".out"
	}

	cmd = strings.Replace(t.cmd, "dotm", cmd, 1) + " --testRoot=" + dir
	out, err := execCommand(cmd)
	if err != nil {
		return fmt.Errorf("Failed to execute '%s'\nOut:%s", cmd, out)
	}
	return nil
}

type file struct {
	path    string
	content string
}

func parseDir(dirname string) ([]*testcase, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	var testcases []*testcase
	for _, file := range files {
		raw, err := ioutil.ReadFile(filepath.Join(dirname, file.Name()))
		if err != nil {
			return nil, err
		}

		c, err := parseTestCase(string(raw))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse %s", file.Name())
		}

		testcases = append(testcases, c)
	}
	return testcases, nil
}

func parseTestCase(raw string) (*testcase, error) {
	sections := strings.Split(string(raw), "---")
	if len(sections) != 4 {
		return nil, fmt.Errorf("expected 4 sections in testcase. got %d", len(sections))
	}
	cmd := strings.TrimSpace(sections[2])
	return &testcase{
		name:     sections[0],
		given:    sections[1],
		cmd:      cmd,
		expected: sections[3],
	}, nil
}
