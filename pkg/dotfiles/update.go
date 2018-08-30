package dotfiles

import (
	"github.com/relnod/dotm/pkg/config"
	"github.com/relnod/dotm/pkg/remote"
)

// Update updates the dotfiles for a given configuration.
func Update(c *config.Config) error {
	err := remote.Pull(c.Remote, c.Repo)
	if err != nil {
		return err
	}

	t := NewTraverser(nil)
	err = t.Traverse(c.Repo, "/tmp/bla2", NewLinkAction(false))
	if err != nil {
		return err
	}

	return nil
}
