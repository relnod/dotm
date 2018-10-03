package remote

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/relnod/fsa"
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
func Clone(fs fsa.FileSystem, remoteURL, path string) error {
	err := fs.MkdirAll(path, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrMkdirDestination)
	}

	remoteURL, path = sanitizePaths(fs, remoteURL, path)

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
func Pull(fs fsa.FileSystem, remoteURL, path string) error {

	remoteURL, path = sanitizePaths(fs, remoteURL, path)
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

func sanitizePaths(fs fsa.FileSystem, remote, path string) (string, string) {
	if filepath.IsAbs(remote) {
		remote, _ = fs.Path(remote)
	}
	path, _ = fs.Path(path)
	return remote, path
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
