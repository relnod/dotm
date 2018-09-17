package remote

import (
	"os"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
)

// Error wrappers
const (
	ErrMkdirDestination = "failed to create destination directory"
	ErrCloneRemote      = "failed to clone remote repository"
	ErrPullRemote       = "failed to pull pull repository"
	ErrOpenLocal        = "failed to open local repository"
)

// Clone clones a remote repository to a given path.
func Clone(remoteURL, path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrMkdirDestination)
	}

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      remoteURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return errors.Wrap(err, ErrCloneRemote)
	}

	return nil
}

// Pull pulles the remote repository.
func Pull(remoteURL, path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return errors.Wrap(err, ErrOpenLocal)
	}

	w, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, ErrOpenLocal)
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrap(err, ErrPullRemote)
	}

	return nil
}

// Detect tries to detect a remote in a given path.
func Detect(path string) (string, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return "", err
	}
	if len(remotes) >= 1 {
		if urls := remotes[0].Config().URLs; len(urls) >= 1 {
			return urls[0], nil
		}
	}

	return "", nil
}
