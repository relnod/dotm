package test

import (
	"os/exec"
)

func execCommand(cmdString string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdString)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
