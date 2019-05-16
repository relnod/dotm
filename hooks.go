package dotm

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

// Hooks represents all hooks.
type Hooks struct {
	PreUpdate  Hook `toml:"pre_update"`
	PostUpdate Hook `toml:"post_update"`
}

// Hook represents one type of hook.
type Hook []string

// ErrExecHook indicates failure during hook exection.
var ErrExecHook = errors.New("failed exec hook")

// Exec executes a hook.
func (h Hook) Exec(ctx context.Context) error {
	for _, cmdRaw := range h {
		args := strings.Split(os.ExpandEnv(cmdRaw), " ")
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("%v: '%s': %s", ErrExecHook, os.ExpandEnv(cmdRaw), err)
		}
	}
	return nil
}

// Merge merges serveral hooks into a single Hooks struct.
func mergeHooks(hooks ...*Hooks) *Hooks {
	merged := &Hooks{}
	for _, h := range hooks {
		merged.PreUpdate = append(merged.PreUpdate, h.PreUpdate...)
		merged.PostUpdate = append(merged.PostUpdate, h.PostUpdate...)
	}
	return merged
}

// ErrOpenHooksFile indicates a failure while opening a hooks file.
var ErrOpenHooksFile = errors.New("failed open hooks file")

// ErrDecodeHooksFile indicates that the hooks file has a syntax error.
var ErrDecodeHooksFile = errors.New("failed decode hooks file")

// LoadHooksFromFile reads a hook file from a given file path.
func LoadHooksFromFile(path string) (*Hooks, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, ErrOpenHooksFile
	}

	var h Hooks
	_, err = toml.Decode(string(data), &h)
	if err != nil {
		return nil, ErrDecodeHooksFile
	}

	return &h, nil
}
