package dotfiles

import (
	"github.com/pkg/errors"

	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Errors
var (
	ErrorEmptyRemote      = errors.New("empty remote url")
	ErrorEmptyDestination = errors.New("empty destination")
)

// Install clones the dotfiles from a remote git repository to a local
// destination and then installs them.
func Install(remoteURL string, destination string) error {
	var err error

	if remoteURL == "" {
		return ErrorEmptyRemote
	}
	if destination == "" {
		return ErrorEmptyDestination
	}

	// @todo: check if dotfiles are already installed

	err = remote.Clone(remoteURL, destination)
	if err != nil {
		return err
	}

	t := NewTraverser(nil)
	err = t.Traverse(destination, "/tmp/bla2", NewLinkAction(false))
	if err != nil {
		return err
	}

	c := &config.Config{
		Remote: remoteURL,
		Repo:   destination,
	}

	err = config.WriteTomlFile("/tmp/bla3/.dotfiles.toml", c)
	if err != nil {
		return err
	}

	return nil
}
