package runner

import "os/exec"

type DotmCmd struct {
	SubCommand string
	Params     map[string]string
}

type Runner struct {
	cmd *exec.Cmd

	err error
}

func Run(cmd DotmCmd) *Runner {
	r := &Runner{}

	r.cmd = exec.Command("dotm")

	r.err = r.cmd.Run()

	return r
}

func (r *Runner) Error() error {
	if r.err != nil {
		return r.err
	}

	// TODO: check for stderr

	return nil
}
