package dotm

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	git "gopkg.in/src-d/go-git.v4"

	"github.com/relnod/dotm/internal/file"
	"github.com/relnod/dotm/internal/fileutil"
)

// Profile defines the data of a dotfile profile.
type Profile struct {
	Name     string   `toml:"-"`
	Path     string   `toml:"path"`
	Remote   string   `toml:"remote"`
	Includes []string `toml:"includes"`
	Excludes []string `toml:"excludes"`
	Hooks

	// expandedPath contains the Path after it was expanded.
	expandedPath string `toml:"-"`
}

// sanitize changes public fields in the profile. Therfore this changes the
// profile.
func (p *Profile) sanitize() {
	p.Remote = sanitizeRemote(p.Remote)
}

// sanitizeRemote checks if the remote is a valid git remote. If not it assumes
// the remote is of the form "domain/user/repo" and converts this to a valid
// https remote.
func sanitizeRemote(remote string) string {
	if strings.HasPrefix(remote, "git@") || strings.HasPrefix(remote, "https://") || remote == "" {
		return remote
	}
	return "https://" + remote + ".git"
}

// expandEnv expands several variables with environment variables.
func (p *Profile) expandEnv() (err error) {
	p.expandedPath, err = expandPath(p.Path)
	if err != nil {
		return err
	}
	return err
}

// expandPath expands the given path with environment variables and converts it
// to an absolute path, if the path is relative.
func expandPath(path string) (string, error) {
	path = os.ExpandEnv(path)

	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	return path, nil
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

// LinkOptions are the options used during the symlink creation.
type LinkOptions struct {
	Force bool
	Dry   bool
	TraversalOptions
}

// link links all files to the home directory.
func (p *Profile) link(opts LinkOptions) error {
	err := p.traverse(&linker{
		source: p.expandedPath,
		dest:   os.Getenv("HOME"),
		force:  opts.Force,
		dry:    opts.Dry,
	}, &opts.TraversalOptions)
	if err != nil {
		return fmt.Errorf("link: %s", err)
	}
	return nil
}

// linker implements fileutil.Visitor
type linker struct {
	// source is the path from where to link from
	source string

	// dest is the path where the files get linked to
	dest string

	force bool
	dry   bool
}

//ErrCreatingDestination indicates failure during the creation of the
//destination dir
var ErrCreatingDestination = errors.New("failed to created destination dir")

func (l *linker) Visit(path, name string) error {
	err := os.MkdirAll(filepath.Join(l.dest, path), os.ModePerm)
	if err != nil {
		return ErrCreatingDestination
	}

	var (
		sourceFile = filepath.Join(l.source, path, name)
		destFile   = filepath.Join(l.dest, path, name)
	)

	// Check if the destination file already exists.
	if _, err := os.Stat(destFile); err == nil {
		if !l.force {
			return nil
		}
		if file.IsSymlink(destFile) {
			err = file.Unlink(destFile, l.dry)
			if err != nil {
				return err
			}
		} else {
			err = file.Backup(destFile, l.dry)
			if err != nil {
				return err
			}
		}
	}

	return file.Link(sourceFile, destFile, l.dry)
}

func (p *Profile) unlink(dry bool) error {
	err := p.traverse(&unlinker{
		path: os.Getenv("HOME"),
		dry:  dry,
	}, nil)
	if err != nil {
		return fmt.Errorf("unlink: %s", err)
	}
	return nil
}

// unlinker implements fileutil.Visitor
type unlinker struct {
	path string

	dry bool
}

func (u *unlinker) Visit(path, name string) error {
	filepath := filepath.Join(u.path, path, name)

	// Check if the file file exists.
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}
	if !file.IsSymlink(filepath) {
		return nil
	}

	err := file.Unlink(filepath, u.dry)
	if err != nil {
		return err
	}

	return file.RestoreBackup(filepath, u.dry)
}

func (p *Profile) addFile(dir, path string) error {
	sourceFile := filepath.Join(os.Getenv("HOME"), path)
	destFile := filepath.Join(p.expandedPath, dir, path)
	err := os.MkdirAll(filepath.Dir(destFile), os.ModePerm)
	if err != nil {
		return err
	}
	data := []byte("# Created by dotm")

	// Check if the source file already exists.
	if _, err := os.Stat(sourceFile); err == nil {
		if file.IsSymlink(sourceFile) {
			err = file.Unlink(sourceFile, false)
			if err != nil {
				return err
			}
		} else {
			data, err = ioutil.ReadFile(sourceFile)
			if err != nil {
				return err
			}
			err = file.Backup(sourceFile, false)
			if err != nil {
				return err
			}
		}
	}

	// Create the file, since it does not exist.
	err = ioutil.WriteFile(destFile, data, os.ModePerm)
	if err != nil {
		return err
	}

	// We link back from the destination file to the source file.
	return file.Link(destFile, sourceFile, false)
}

// FindHooks finds all hooks of a given profile.
// Hooks can be found at:
// - ~/.config/dotm/config.toml
// - ~/.config/dotm/profiles/<profile>/hooks.toml
// - ~/.config/dotm/profiles/<profile>/<top-level-dir>/hooks.toml
func (p *Profile) findHooks(opts *TraversalOptions) (*Hooks, error) {
	var hooks []*Hooks

	hooks = append(hooks, &p.Hooks)

	h, err := findHook(p.expandedPath)
	if err != nil {
		return nil, err
	}
	if h != nil {
		hooks = append(hooks, h)
	}

	topLevelDirs, err := p.topLevelDirs(opts)
	if err != nil {
		return nil, err
	}
	for _, dir := range topLevelDirs {
		h, err := findHook(filepath.Join(p.expandedPath, dir))
		if err != nil {
			return nil, err
		}
		if h != nil {
			hooks = append(hooks, h)
		}
	}

	return mergeHooks(hooks...), nil
}

const hooksFileName = "hooks.toml"

func findHook(dir string) (*Hooks, error) {
	filepath := filepath.Join(dir, hooksFileName)
	if _, err := os.Stat(filepath); err != nil {
		return nil, nil
	}

	return LoadHooksFromFile(filepath)
}

var (
	// ErrCreateLocalPath indicates a failure during the creation of the local
	// path.
	ErrCreateLocalPath = errors.New("failed to create local path")

	// ErrCloneRemote indicates an unsuccesful remote git clone.
	ErrCloneRemote = errors.New("failed to clone remote")
)

// cloneRemote clones the remote repository to the local path.
func (p *Profile) cloneRemote(ctx context.Context) error {
	err := os.MkdirAll(p.expandedPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("clone: %v", ErrCreateLocalPath)
	}

	_, err = git.PlainCloneContext(ctx, p.expandedPath, false, &git.CloneOptions{
		URL: p.Remote,
	})
	if err != nil {
		return fmt.Errorf("clone: %v: %v", ErrCloneRemote, err)
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
func (p *Profile) pullRemote(ctx context.Context) error {
	r, err := git.PlainOpen(p.expandedPath)
	if err != nil {
		return ErrOpenRepo
	}

	w, err := r.Worktree()
	if err != nil {
		return ErrOpenRepo
	}

	err = w.PullContext(ctx, &git.PullOptions{RemoteName: "origin"})
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

// TraversalOptions are used during the dotfile traversal.
type TraversalOptions struct {
	Includes []string
	Excludes []string
}

func (p *Profile) traverse(visitor fileutil.Visitor, opts *TraversalOptions) error {
	topLevelDirs, err := p.topLevelDirs(opts)
	if err != nil {
		return err
	}

	for _, d := range topLevelDirs {
		err := fileutil.RecTraverseDir(filepath.Join(p.expandedPath, d), visitor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) topLevelDirs(opts *TraversalOptions) ([]string, error) {
	files, err := ioutil.ReadDir(p.expandedPath)
	if err != nil {
		return nil, err
	}

	includes := p.Includes
	excludes := p.Excludes
	if opts != nil {
		includes = append(includes, opts.Includes...)
		excludes = append(excludes, opts.Excludes...)
	}

	dirs := []string{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if !isIncluded(f.Name(), includes) {
			continue
		}
		if isExcluded(f.Name(), excludes) {
			continue
		}

		dirs = append(dirs, f.Name())
	}

	return dirs, nil
}

// alwaysExcluded is a list of directories, that always get excluded.
var alwaysExcluded = []string{".git"}

// isExcluded checks if the dir should be excluded.
// Also excludes all directories prefixed with a "_".
func isExcluded(dir string, excludes []string) bool {
	if strings.HasPrefix(dir, "_") {
		return true
	}
	excludes = append(excludes, alwaysExcluded...)
	for _, exclude := range excludes {
		if dir == exclude {
			return true
		}
	}
	return false
}

// isIncluded checks if the directory should be included.
// When the includes list is empty, it returns true.
func isIncluded(dir string, includes []string) bool {
	if len(includes) == 0 {
		return true
	}
	for _, include := range includes {
		if dir == include {
			return true
		}
	}
	return false
}
