package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

type testcase struct {
	name     string
	given    string
	cmd      string
	expected string
}

func (t *testcase) exec(dir string) error {
	args := strings.Split(t.cmd, " ")
	args = append(args, fmt.Sprintf("--testRoot=%s", dir))
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute '%s'\nStdout:\n%s\nStderr:\n%s", t.cmd, out.String(), outErr.String())
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
			return nil, err
		}

		testcases = append(testcases, c)
	}
	return testcases, nil
}

func parseTestCase(raw string) (*testcase, error) {
	sections := strings.Split(string(raw), "---")
	if len(sections) != 4 {
		return nil, fmt.Errorf("Expected 4 sections in testcase. got %d", len(sections))
	}
	return &testcase{
		name:     sections[0],
		given:    sections[1],
		cmd:      strings.Replace(sections[2], "\n", "", -1),
		expected: sections[3],
	}, nil
}
