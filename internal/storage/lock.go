package storage

import (
	"os"
	"strconv"
)

func IsLockFileExists(cacheRoot string) bool {
	lockPath := cacheRoot + dotLock
	_, err := os.Stat(lockPath)
	return err == nil
}

func CreateLockFile(cacheRoot string) error {
	lockPath := cacheRoot + dotLock
	pid := os.Getpid()
	return os.WriteFile(lockPath, []byte(strconv.Itoa(pid)), 0644)
}

func RemoveLockFile(cacheRoot string) error {
	lockPath := cacheRoot + dotLock
	return os.Remove(lockPath)
}
