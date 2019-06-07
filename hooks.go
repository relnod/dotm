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
	PreUpdate  Hook `toml:"pre_update" clic:"pre_update"`
	PostUpdate Hook `toml:"post_update" clic:"post_update"`
}

// Hook represents one type of hook.
type Hook []string

// Exec executes a hook.
func (h Hook) Exec(ctx context.Context) error {
	for _, cmdRaw := range h {
		args := strings.Split(os.ExpandEnv(cmdRaw), " ")
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed exec hook: '%s': %s", os.ExpandEnv(cmdRaw), err)
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

// LoadHooksFromFile reads a hook file from a given file path.
func LoadHooksFromFile(path string) (*Hooks, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed open hooks file: %v", err)
	}

	var h Hooks
	_, err = toml.Decode(string(data), &h)
	if err != nil {
		return nil, errors.New("failed decode hooks file")
	}

	return &h, nil
}
