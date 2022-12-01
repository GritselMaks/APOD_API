package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func InitFile(path string, fileMaxSizeBytes int64) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
}

func StableFilePath(path string) (*string, error) {
	return preparePath(path)
}

func preparePath(path string) (*string, error) {
	if len(path) == 0 || (path[0] != '~' && path[0] != '/') {
		return nil, fmt.Errorf("invalid input")
	}
	if path[0] == '/' {
		return &path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(usr.HomeDir, path[1:])
	return &path, nil
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--
			continue
		}

		return nil
	}

	return
}
