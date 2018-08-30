package runner

import "os/exec"

// DotmCmd wraps a dotm command.
type DotmCmd struct {
	SubCommand string
	Params     map[string]string
}

func (d DotmCmd) buildParams() []string {
	p := make([]string, len(d.Params))

	p = append(p, d.SubCommand)

	for arg, val := range d.Params {
		p = append(p, "--"+arg+"="+val)
	}

	return p
}

// Runner runs a dotm command and provides helper methods to process output.
type Runner struct {
	cmd *exec.Cmd

	err error
}

// Run runs a dotm command.
func Run(cmd DotmCmd) *Runner {
	r := &Runner{}

	r.cmd = exec.Command("dotm", cmd.buildParams()...)

	r.err = r.cmd.Run()

	return r
}

// Error returns an error if the dotm command failed.
func (r *Runner) Error() error {
	if r.err != nil {
		return r.err
	}

	// TODO: check for stderr

	return nil
}
