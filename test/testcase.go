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
	name      string
	given     string
	cmd       string
	cmdOutput string
	expected  string
}

func (t *testcase) exec(cmd, dir string, index int) (string, error) {
	if coverage {
		cmd += " -test.coverprofile=$ROOT/coverage-e2e-" + strconv.Itoa(index) + ".out"
	}

	cmd = strings.Replace(t.cmd, "dotm", cmd, 1) + " --testRoot=" + dir
	out, err := execCommand(cmd)
	if err != nil {
		return out, fmt.Errorf("failed to execute '%s'\nOut:%s", cmd, out)
	}
	return out, err
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

		rawCases := strings.Split(string(raw), "===")
		for _, rawCase := range rawCases {
			c, err := parseTestCase(rawCase)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse %s", file.Name())
			}
			testcases = append(testcases, c)
		}
	}
	return testcases, nil
}

func parseTestCase(raw string) (*testcase, error) {
	sections := strings.Split(string(raw), "---")
	if len(sections) != 4 {
		return nil, fmt.Errorf("expected 4 sections in testcase. got %d", len(sections))
	}
	cmdSplit := strings.Split(strings.TrimSpace(sections[2]), ":")
	cmd := cmdSplit[0]
	cmdOutput := ""
	if len(cmdSplit) == 2 {
		cmdOutput = strings.Replace(cmdSplit[1], "\\n", "\n", -1)
	}
	return &testcase{
		name:      sections[0][:len(sections[0])-1],
		given:     sections[1],
		cmd:       cmd,
		cmdOutput: cmdOutput,
		expected:  sections[3],
	}, nil
}
