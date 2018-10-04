package remote

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/relnod/dotm/pkg/config"
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

// CloneProfile clones the remote repository for a profile.
func CloneProfile(fs fsa.FileSystem, p *config.Profile) error {
	err := fs.MkdirAll(p.Path, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrMkdirDestination)
	}

	remoteURL, path := sanitizePaths(fs, p.Remote, p.Path)

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      remoteURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return errors.Wrap(err, ErrCloneRemote)
	}

	return nil
}

// PullProfile pulles the remote repository for a profile
func PullProfile(fs fsa.FileSystem, p *config.Profile) error {
	_, path := sanitizePaths(fs, p.Remote, p.Path)
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
