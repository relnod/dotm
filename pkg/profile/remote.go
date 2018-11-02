package profile

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

// CloneRemote clones the remote repository for a profile..
func (p *Profile) CloneRemote() error {
	remoteURL, path := sanitizePaths(p.fs, p.Remote, p.Path)

	err := p.fs.MkdirAll(p.Path, os.ModePerm)
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

// PullRemote pulles the remote repository for a profile.
func (p *Profile) PullRemote() error {
	_, path := sanitizePaths(p.fs, p.Remote, p.Path)

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

// DetectRemote tries to detect the remote for a profile.
func (p *Profile) DetectRemote() (string, error) {
	_, path := sanitizePaths(p.fs, p.Remote, p.Path)

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

func sanitizePaths(fs fsa.FileSystem, remote, path string) (string, string) {
	if filepath.IsAbs(remote) {
		remote, _ = fs.Path(remote)
	}
	path, _ = fs.Path(path)
	return remote, path
}
