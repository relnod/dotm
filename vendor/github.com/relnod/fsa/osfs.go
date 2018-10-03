package fsa

import (
	"os"
	"time"
)

type OsFs struct{}

func NewOsFs() *OsFs {
	return &OsFs{}
}

func (f *OsFs) Path(path string) (string, error) {
	return path, nil
}

func (f *OsFs) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (f *OsFs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (f *OsFs) Remove(name string) error {
	return os.Remove(name)
}

func (f *OsFs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (f *OsFs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (f *OsFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (f *OsFs) Chown(name string, uid, gid int) error {
	return os.Chown(name, uid, gid)
}

func (f *OsFs) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (f *OsFs) Getwd() (dir string, err error) {
	return os.Getwd()
}

func (f *OsFs) TempDir() string {
	return os.TempDir()
}

func (f *OsFs) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (f *OsFs) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f *OsFs) MkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}

func (f *OsFs) RemoveAll(path string) (err error) {
	return os.RemoveAll(path)
}

func (f *OsFs) Truncate(name string, size int64) error {
	return os.Truncate(name, size)
}

func (f *OsFs) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (f *OsFs) Lchown(name string, uid, gid int) error {
	return os.Lchown(name, uid, gid)
}

func (f *OsFs) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

func (f *OsFs) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}
