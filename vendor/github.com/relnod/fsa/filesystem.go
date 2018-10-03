package fsa

import (
	"os"
	"time"
)

type FileSystem interface {
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	Mkdir(name string, perm os.FileMode) error
	Remove(name string) error
	Stat(name string) (os.FileInfo, error)
	Chmod(name string, mode os.FileMode) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
	Chown(name string, uid, gid int) error

	Chdir(dir string) error
	Getwd() (dir string, err error)
	TempDir() string
	Open(name string) (*os.File, error)
	Create(name string) (*os.File, error)
	MkdirAll(name string, perm os.FileMode) error
	RemoveAll(path string) (err error)
	Truncate(name string, size int64) error

	Lstat(name string) (os.FileInfo, error)
	Lchown(name string, uid, gid int) error
	Readlink(name string) (string, error)
	Symlink(oldname, newname string) error

	Path(string) (string, error)
}
