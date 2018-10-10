package hook

import (
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"
)

// Errors
var (
	ErrOpenHooksFile   = "failed to open hooks file '%s'"
	ErrDecodeHooksFile = "failed to decode hooks file '%s'"
	ErrExecuteHook     = "failed to execute hook '%s'"
)

// Hooks represents all hooks.
type Hooks struct {
	PreUpdate  Hook `toml:"pre_update"`
	PostUpdate Hook `toml:"post_update"`
}

// Hook represents one type of hook.
type Hook []string

// Execute executes a hook.
func (h Hook) Execute() error {
	for _, cmdRaw := range h {
		args := strings.Split(os.ExpandEnv(cmdRaw), " ")
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return errors.Wrapf(err, ErrExecuteHook, cmdRaw)
		}
	}
	return nil
}

// Merge merges serveral hooks into a single Hooks struct.
func Merge(hooks ...*Hooks) *Hooks {
	merged := &Hooks{}
	for _, h := range hooks {
		merged.PreUpdate = append(merged.PreUpdate, h.PreUpdate...)
		merged.PostUpdate = append(merged.PostUpdate, h.PostUpdate...)
	}
	return merged
}

// ReadFromFile reads a new hook from a given file path.
func ReadFromFile(fs fsa.FileSystem, path string) (*Hooks, error) {
	data, err := fsutil.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrapf(err, ErrOpenHooksFile, path)
	}

	var h Hooks
	_, err = toml.Decode(string(data), &h)
	if err != nil {
		return nil, errors.Wrapf(err, ErrDecodeHooksFile, path)
	}

	return &h, nil
}
