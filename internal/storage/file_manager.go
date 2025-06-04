package storage

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

var FileManager *fileManager

type fileManager struct {
	targetUser *user.User
	uid        int
	gid        int
}

func init() {
	var targetUser *user.User
	var gid int
	var uid int

	if sudoUser := os.Getenv(sudoUserEnv); sudoUser != "" {
		var err error
		if targetUser, err = user.Lookup(sudoUser); err == nil {
			uid, _ = strconv.Atoi(targetUser.Uid)
			gid, _ = strconv.Atoi(targetUser.Gid)
		}
	}

	FileManager = &fileManager{
		targetUser: targetUser,
		uid:        uid,
		gid:        gid,
	}
}

func (fm *fileManager) WriteFile(path string, data []byte, mode os.FileMode) error {
	if err := os.WriteFile(path, data, mode); err != nil {
		return err
	}

	if fm.targetUser != nil {
		return fm.Chown(path)
	}

	return nil
}

func (fm *fileManager) Remove(path string) error {
	return os.Remove(path)
}

func (fm *fileManager) MkdirAll(dir string, mode os.FileMode) error {
	var parentDir string
	var needsParentChown bool
	if fm.targetUser != nil {
		parentDir = filepath.Dir(dir)
		_, err := os.Stat(parentDir)
		needsParentChown = err != nil
	}

	if err := os.MkdirAll(dir, mode); err != nil {
		return err
	}

	if fm.targetUser != nil {
		if needsParentChown {
			if err := fm.Chown(parentDir); err != nil {
				return err
			}
		}

		return fm.Chown(dir)
	}

	return nil
}

func (fm *fileManager) Chown(dir string) error {
	return os.Chown(dir, fm.uid, fm.gid)
}
