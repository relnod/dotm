package fileutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"
	"github.com/relnod/fsa/testutil"
)

// Errors
var (
	ErrCreatingSymlink = errors.New("Failed to create Symlink")
)

// Visitor defines a visitor.
type Visitor interface {
	// Visit gets called for each file being traversed.
	Visit(dir string, file string)
}

// RecTraverseDir recursively traverses all directories starting at dir.
// Calls the visitor for each file it passes and passes the relDir and file name
// to the visitor.
func RecTraverseDir(fs fsa.FileSystem, dir string, relDir string, visitor Visitor) error {
	files, err := fsutil.ReadDir(fs, dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			err := RecTraverseDir(fs, filepath.Join(dir, f.Name()), filepath.Join(relDir, f.Name()), visitor)
			if err != nil {
				return err
			}
		} else {
			visitor.Visit(relDir, f.Name())
		}
	}

	return nil
}

// Link creates a symbolik link from dest to source. When dry is true only
// perfomers a dry run.
func Link(fs fsa.FileSystem, from string, to string, dry bool) error {
	if dry {
		fmt.Printf("Creating Symlink from %s to %s\n", from, to)
		return nil
	}

	err := fs.Symlink(from, to)
	if err != nil {
		return ErrCreatingSymlink
	}

	return nil
}

// Unlink removes a symbolik link from dest to source. When dry is true only
// perfomers a dry run.
func Unlink(fs fsa.FileSystem, file string, dry bool) error {
	if dry {
		fmt.Printf("Removing %s\n", file)
		return nil
	}
	err := fs.Remove(file)
	if err != nil {
		return err
	}

	return nil
}

// Backup moves the given file to filename.backup. When dry is true only
// perfomers a dry run.
func Backup(fs fsa.FileSystem, file string, dry bool) error {
	if dry {
		fmt.Printf("Creating backup for %s\n", file)
		return nil
	}
	return moveFile(fs, file, backupPath(file))
}

// RestoreBackup tries to restore a backup. When dry is true only
// perfomers a dry run.
func RestoreBackup(fs fsa.FileSystem, file string, dry bool) error {
	if !testutil.FileExists(fs, backupPath(file)) {
		spew.Dump(backupPath(file))
		return nil
	}
	if dry {
		fmt.Printf("Creating backup for %s\n", file)
		return nil
	}

	return moveFile(fs, backupPath(file), file)
}

func moveFile(fs fsa.FileSystem, from, to string) error {
	data, err := fsutil.ReadFile(fs, from)
	if err != nil {
		return err
	}
	err = fsutil.WriteFile(fs, to, data, os.ModePerm)
	if err != nil {
		return err
	}
	return fs.Remove(from)
}

func backupPath(path string) string {
	return path + ".backup"
}
