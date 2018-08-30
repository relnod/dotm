package dotfiles

import (
	"os"

	"github.com/pkg/errors"
	"github.com/relnod/dotm/pkg/config"
	git "gopkg.in/src-d/go-git.v4"
)

// Errors
var (
	ErrorEmptyRemote      = errors.New("empty remote url")
	ErrorEmptyDestination = errors.New("empty destination")
)

// Error Wrappings
const (
	ErrorMkdirDestination = "failed to create destination directory"
	ErrorCloneRemote      = "failed to clone remote repository"
)

// Install clones the dotfiles from a remote git repository to a local
// destination and then installs them.
func Install(remote string, destination string) error {
	var err error

	if remote == "" {
		return ErrorEmptyRemote
	}
	if remote == "" {
		return ErrorEmptyRemote
	}

	// @todo: check if dotfiles are already installed

	err = cloneRemote(remote, destination)
	if err != nil {
		return err
	}

	t := NewTraverser(nil)
	err = t.Traverse(destination, "/tmp/bla2", NewLinkAction(false))
	if err != nil {
		return err
	}

	c := &config.Config{
		Remote: remote,
		Repo:   destination,
	}

	err = config.WriteTomlFile("/tmp/bla3/.dotfiles.toml", c)
	if err != nil {
		return err
	}

	return nil
}

func cloneRemote(remote, destination string) error {
	err := os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrorMkdirDestination)
	}

	_, err = git.PlainClone(destination, false, &git.CloneOptions{
		URL:      remote,
		Progress: os.Stdout,
	})
	if err != nil {
		errors.Wrap(err, ErrorCloneRemote)
	}

	return nil
}
