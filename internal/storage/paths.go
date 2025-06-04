package storage

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func GetCachePath() (string, error) {
	userCacheDir, err := GetUserCachePath()
	if err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	cachePath := filepath.Join(userCacheDir, qpCacheDir)
	if err := FileManager.MkdirAll(cachePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cachePath, nil
}

func GetUserCachePath() (string, error) {
	err := switchToRealUser()
	if err != nil {
		return "", err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv(homeEnv)
	}

	if runtime.GOOS == osDarwin {
		return filepath.Join(home, darwinCacheDir), nil
	}

	userCacheDir := os.Getenv(xdgCacheHomeEnv)
	if userCacheDir == "" {
		userCacheDir = filepath.Join(home, dotCache)
	}

	return userCacheDir, nil
}

func switchToRealUser() error {
	realUser := os.Getenv(sudoUserEnv)
	if realUser == "" {
		return nil
	}

	usr, err := user.Lookup(realUser)
	if err != nil {
		return err
	}

	os.Setenv(homeEnv, usr.HomeDir)
	os.Setenv(userEnv, realUser)
	os.Unsetenv(xdgCacheHomeEnv)

	return nil
}
