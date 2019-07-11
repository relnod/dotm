package dotm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	git "gopkg.in/src-d/go-git.v4"

	"github.com/relnod/dotm/internal/file"
	"github.com/relnod/dotm/internal/fileutil"
)

// Profile defines the data of a dotfile profile.
type Profile struct {
	Name             string            `toml:"-" clic:"name"`
	Path             string            `toml:"path" clic:"path"`
	Remote           string            `toml:"remote" clic:"remote"`
	HooksEnabled     bool              `toml:"hooks_enabled" clic:"hooks_enabled"`
	UpdateFromRemote bool              `toml:"update_from_remote" clic:"update_from_remote"`
	Includes         []string          `toml:"includes" clic:"includes"`
	Excludes         []string          `toml:"excludes" clic:"excludes"`
	Vars             map[string]string `toml:"vars" clic:"vars"`
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

// Create creates the path of the profile.
// It also initializes a new git repository in the profile root.
func (p *Profile) create() error {
	err := os.MkdirAll(p.expandedPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create: %v", errCreateLocalPath)
	}

	_, err = git.PlainInit(p.expandedPath, false)
	if err != nil {
		return fmt.Errorf("create: failed to initialize git repo: %v", err)
	}
	return nil
}

// LinkOptions are the options used during the symlink creation.
type LinkOptions struct {
	Force bool
	Dry   bool
	TraversalOptions
}

// link symlinks all profile files to the home directory.
//
// Examples:
//      <profilepath>/bash/.bashrc -> $HOME/.bashrc
//      <profilepath>/nvim/.config/nvim/init.vim -> $HOME/.config/nvim/init.vim
func (p *Profile) link(opts LinkOptions) error {
	err := p.traverse(&linker{
		dest:    os.Getenv("HOME"),
		force:   opts.Force,
		dry:     opts.Dry,
		tplVars: p.Vars,
	}, &opts.TraversalOptions)
	if err != nil {
		return fmt.Errorf("link: %s", err)
	}
	return nil
}

// linker implements fileutil.Visitor
type linker struct {
	// dest is the path where the files get linked to
	dest string

	force   bool
	dry     bool
	tplVars map[string]string
}

func (l *linker) Visit(path, relativePath string) error {
	err := os.MkdirAll(filepath.Join(l.dest, filepath.Dir(relativePath)), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to created destination dir: %v", err)
	}

	destFile := filepath.Join(l.dest, relativePath)

	// Don't link already processed templates, since they should have been
	// handled below.
	if strings.HasSuffix(path, ".tpl.out") {
		return nil
	}
	// Check if the file is a template
	if strings.HasSuffix(path, ".tpl") {
		path, err = processTemplate(path, l.tplVars)
		if err != nil {
			return err
		}
		// Remove the .tpl suffix of the destination file.
		destFile = strings.TrimSuffix(destFile, ".tpl")
	}

	// Check if the file is a symlink. If so remove it, even if the force option
	// is not set.
	if file.IsSymlink(destFile) {
		err = file.Unlink(destFile, l.dry)
		if err != nil {
			return err
		}

		// Check if the destination file already exists.
		// If it exists create a backup file.
	} else if _, err := os.Stat(destFile); err == nil {
		if !l.force {
			return nil
		}
		err = file.Backup(destFile, l.dry)
		if err != nil {
			return err
		}
	}

	return file.Link(path, destFile, l.dry)
}

// processTemplate generates the processed template file from the given file
// path. All vars are used for tempalte parsing.
// It returns the path of the generated file.
// It returns an error, when parsing fails or a filesystem problem occurs.
func processTemplate(path string, vars map[string]string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("tpl: %v", err)
	}
	tmpl, err := template.New(path).Parse(string(data))
	if err != nil {
		return "", fmt.Errorf("tpl: %v", err)
	}
	b := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(b, vars)
	if err != nil {
		return "", fmt.Errorf("tpl: %v", err)
	}
	path = path + ".out"
	err = ioutil.WriteFile(path, b.Bytes(), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("tpl: %v", err)
	}
	return path, nil
}

// unlink removes all symlinks created by the profile.
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

func (u *unlinker) Visit(path, relativePath string) error {
	filepath := filepath.Join(u.path, relativePath)

	// Skip template files.
	if strings.HasSuffix(filepath, ".tpl") {
		return nil
	}
	// Remove suffix from processed template files.
	if strings.HasSuffix(filepath, ".tpl.out") {
		filepath = strings.TrimSuffix(filepath, ".tpl.out")
	}

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

// addFile adds a file to the profile. The file gets added to the given top
// level directory. If the file already exists under $HOME/path, the file gets
// moved to its new location in the profile. Otherwise a new file gets created.
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

	h, err := findAndLoadHook(p.expandedPath)
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
		h, err := findAndLoadHook(filepath.Join(p.expandedPath, dir))
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

// findAndLoadHook checks if a hooks file exists at the given directory.
// If a hooks file exists, it gets loaded.
func findAndLoadHook(dir string) (*Hooks, error) {
	filepath := filepath.Join(dir, hooksFileName)
	if _, err := os.Stat(filepath); err != nil {
		return nil, nil
	}

	return LoadHooksFromFile(filepath)
}

// errCreateLocalPath indicates a failure during the creation of the local
// path.
var errCreateLocalPath = errors.New("failed to create local path")

// cloneRemote clones the remote git repository to the local path.
func (p *Profile) cloneRemote(ctx context.Context) error {
	err := os.MkdirAll(p.expandedPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("clone: %v", errCreateLocalPath)
	}

	_, err = git.PlainCloneContext(ctx, p.expandedPath, false, &git.CloneOptions{
		URL: p.Remote,
	})
	if err != nil {
		return fmt.Errorf("clone: failed to clone remote: %v", err)
	}
	return nil
}

// errOpenRepo indicates a failure while opening a git repository.
var errOpenRepo = errors.New("failed to open repository")

// pullRemote pulls updates from the remote git repository.
func (p *Profile) pullRemote(ctx context.Context) error {
	r, err := git.PlainOpen(p.expandedPath)
	if err != nil {
		return fmt.Errorf("pull: %v: %v", errOpenRepo, err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("pull: %v: %v", errOpenRepo, err)
	}

	status, err := w.Status()
	if err != nil {
		return err
	}
	if !status.IsClean() {
		return errors.New("pull: worktree is not clean")
	}

	err = w.PullContext(ctx, &git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("pull: failed to pull remote: %v", err)
	}
	return nil
}

// detectRemote tries to detect the git remote path.
func (p *Profile) detectRemote() (string, error) {
	r, err := git.PlainOpen(p.expandedPath)
	if err != nil {
		return "", errOpenRepo
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
	Includes     []string
	Excludes     []string
	IgnorePrefix string
}

// traverse traverses the files of the profile.
func (p *Profile) traverse(visitor fileutil.Visitor, opts *TraversalOptions) error {
	if opts == nil {
		opts = &TraversalOptions{}
	}

	topLevelDirs, err := p.topLevelDirs(opts)
	if err != nil {
		return err
	}

	for _, d := range topLevelDirs {
		err := fileutil.RecTraverseDir(filepath.Join(p.expandedPath, d), visitor, opts.IgnorePrefix)
		if err != nil {
			return err
		}
	}
	return nil
}

// topLevelDirs returns all top level directories of the profile.
//
// Top level directories can be limitd by the profile includes/excludes or by
// specifying them in the TraversalOptions. Includes and excludes from both
// sources get merged.
//
// Includes take precedence over excludes.
func (p *Profile) topLevelDirs(opts *TraversalOptions) ([]string, error) {
	if opts == nil {
		opts = &TraversalOptions{}
	}

	files, err := ioutil.ReadDir(p.expandedPath)
	if err != nil {
		return nil, err
	}

	includes := append(p.Includes, opts.Includes...)
	excludes := append(p.Excludes, opts.Excludes...)

	dirs := []string{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if !isIncluded(f.Name(), includes) {
			continue
		}
		if isExcluded(f.Name(), excludes, opts.IgnorePrefix) {
			continue
		}

		dirs = append(dirs, f.Name())
	}

	return dirs, nil
}

// alwaysExcluded is a list of directories, that always get excluded.
var alwaysExcluded = []string{".git"}

// isExcluded checks if the dir should be excluded.
// Also excludes all directories prefixed with the ignorePrefix.
func isExcluded(dir string, excludes []string, ignorePrefix string) bool {
	if ignorePrefix != "" && strings.HasPrefix(dir, ignorePrefix) {
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
