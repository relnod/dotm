package remote

import (
	"os"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
)

// Error wrappers
const (
	ErrorMkdirDestination = "failed to create destination directory"
	ErrorCloneRemote      = "failed to clone remote repository"
	ErrorPullRemote       = "failed to pull pull repository"
	ErrorOpenLocal        = "failed to open local repository"
)

// Clone clones a remote repository to a given path.
func Clone(remoteURL, path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrorMkdirDestination)
	}

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      remoteURL,
		Progress: os.Stdout,
	})
	if err != nil {
		errors.Wrap(err, ErrorCloneRemote)
	}

	return nil
}

// Pull pulles the remote repository.
func Pull(remoteURL, path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return errors.Wrap(err, ErrorOpenLocal)
	}

	w, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, ErrorOpenLocal)
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return errors.Wrap(err, ErrorPullRemote)
	}

	return nil
}
