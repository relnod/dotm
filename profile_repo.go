package dotm

import (
	"errors"
	"os"

	"golang.org/x/xerrors"
	git "gopkg.in/src-d/go-git.v4"
)

var (
	// ErrCreateLocalPath indicates a failure during the creation of the local
	// path.
	ErrCreateLocalPath = errors.New("failed to create local path")

	// ErrCloneRemote indicates an unsuccesful remote git clone.
	ErrCloneRemote = errors.New("failed to clone remote")
)

// cloneRemote clones the remote repository to the local path.
func (p *Profile) cloneRemote() error {
	err := os.MkdirAll(p.expandedPath, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("clone: %v", ErrCreateLocalPath)
	}

	_, err = git.PlainClone(p.expandedPath, false, &git.CloneOptions{
		URL: p.Remote,
	})
	if err != nil {
		return xerrors.Errorf("clone: %v: %v", ErrCloneRemote, err)
	}
	return nil
}

var (
	// ErrOpenRepo indicates a failure while opening a git repository.
	ErrOpenRepo = errors.New("failed to open repository")

	// ErrPullRemote indicates an unsuccesful remote git pull.
	ErrPullRemote = errors.New("failed to pull remote")
)

// pullRemote pulls updates from the remote repository.
func (p *Profile) pullRemote() error {
	r, err := git.PlainOpen(p.expandedPath)
	if err != nil {
		return ErrOpenRepo
	}

	w, err := r.Worktree()
	if err != nil {
		return ErrOpenRepo
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return ErrPullRemote
	}
	return nil
}

func (p *Profile) detectRemote() (string, error) {
	r, err := git.PlainOpen(p.expandedPath)
	if err != nil {
		return "", ErrOpenRepo
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

// ErrInitRepo indicates an unsuccesful git init
var ErrInitRepo = errors.New("failed to initialize git repo")

// Create creates the path of the profile.
func (p *Profile) create() error {
	err := os.MkdirAll(p.expandedPath, os.ModePerm)
	if err != nil {
		return ErrCreateLocalPath
	}

	_, err = git.PlainInit(p.expandedPath, false)
	if err != nil {
		return ErrInitRepo
	}
	return nil
}
