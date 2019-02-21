package util

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"os"
)

func FileExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

func DirExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}


func HomeDir() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return dir, nil
}

func MakeDirIfNotExists(path string, perm os.FileMode) error {
	if !DirExists(path) {
		if err := os.MkdirAll(path, perm); err != nil {
			return errors.Wrapf(err, "failed to create directory %v", err)
		}
	}
	return nil
}
