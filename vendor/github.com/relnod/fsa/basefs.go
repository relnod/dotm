package fsa

import (
	"os"
	"path/filepath"
	"time"
)

type BaseFs struct {
	f    FileSystem
	base string
}

func NewBaseFs(fs FileSystem, base string) *BaseFs {
	fs.MkdirAll(base, os.ModePerm)
	return &BaseFs{fs, base}
}

func (f *BaseFs) Cleanup() error {
	return f.f.RemoveAll(f.base)
}

func (f *BaseFs) Base() string {
	return f.base
}
func (f *BaseFs) path(path string) (string, error) {
	if !filepath.IsAbs(path) {
		wd, err := f.f.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(wd, path)

	}
	return filepath.Join(f.base, path), nil
}

func (f *BaseFs) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	name2, err := f.path(name)
	if err != nil {
		return nil, err
	}
	return f.f.OpenFile(name2, flag, perm)
}

func (f *BaseFs) Mkdir(name string, perm os.FileMode) error {
	name2, err := f.path(name)
	if err != nil {
		return err
	}
	return f.f.Mkdir(name2, perm)
}

func (f *BaseFs) Remove(name string) error {
	name2, err := f.path(name)
	if err != nil {
		return err
	}
	return f.f.Remove(name2)
}

func (f *BaseFs) Stat(name string) (os.FileInfo, error) {
	name2, err := f.path(name)
	if err != nil {
		return nil, err
	}
	return f.f.Stat(name2)
}

func (f *BaseFs) Chmod(name string, mode os.FileMode) error {
	name2, err := f.path(name)
	if err != nil {
		return err
	}
	return f.f.Chmod(name2, mode)
}

func (f *BaseFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	name2, err := f.path(name)
	if err != nil {
		return err
	}
	return f.f.Chtimes(name2, atime, mtime)
}

func (f *BaseFs) Chown(name string, uid, gid int) error {
	name2, err := f.path(name)
	if err != nil {
		return err
	}
	return f.f.Chown(name2, uid, gid)
}

func (f *BaseFs) Chdir(dir string) error {
	dir2, err := f.path(dir)
	if err != nil {
		return err
	}
	return f.f.Chdir(dir2)
}

func (f *BaseFs) Getwd() (dir string, err error) {
	wd, err := f.f.Getwd()
	if err != nil {
		return "", err
	}
	wd2, err := f.path(wd)
	if err != nil {
		return "", err
	}
	return wd2, nil
}

func (f *BaseFs) TempDir() string {
	tempDir := f.f.TempDir()
	tempDir2, err := f.path(tempDir)
	if err != nil {

	}
	return tempDir2
}

func (f *BaseFs) Open(name string) (*os.File, error) {
	name2, err := f.path(name)
	if err != nil {

	}
	return f.f.Open(name2)
}

func (f *BaseFs) Create(name string) (*os.File, error) {
	name2, err := f.path(name)
	if err != nil {

	}
	return f.f.Create(name2)
}

func (f *BaseFs) MkdirAll(name string, perm os.FileMode) error {
	name2, err := f.path(name)
	if err != nil {

	}
	return f.f.MkdirAll(name2, perm)
}

func (f *BaseFs) RemoveAll(path string) (err error) {
	path2, err := f.path(path)
	if err != nil {

	}
	return f.f.RemoveAll(path2)
}

func (f *BaseFs) Truncate(name string, size int64) error {
	name2, err := f.path(name)
	if err != nil {

	}
	return f.f.Truncate(name2, size)
}

func (f *BaseFs) Lstat(name string) (os.FileInfo, error) {
	name2, err := f.path(name)
	if err != nil {

	}
	return f.f.Lstat(name2)
}

func (f *BaseFs) Lchown(name string, uid, gid int) error {
	name2, err := f.path(name)
	if err != nil {

	}
	return f.f.Lchown(name2, uid, gid)
}

func (f *BaseFs) Readlink(name string) (string, error) {
	name2, err := f.path(name)
	if err != nil {

	}
	return f.f.Readlink(name2)
}

func (f *BaseFs) Symlink(oldname, newname string) error {
	oldname2, err := f.path(oldname)
	if err != nil {

	}
	newname2, err := f.path(newname)
	if err != nil {

	}
	return f.f.Symlink(oldname2, newname2)
}
