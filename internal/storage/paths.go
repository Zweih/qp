package storage

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"qp/internal/consts"
	"runtime"
)

func GetCachePath() (string, error) {
	userCacheDir, err := GetUserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	cacheDir := filepath.Join(userCacheDir, qpCacheDir)
	if err := FileManager.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}

func GetUserCacheDir() (string, error) {
	err := switchToRealUser()
	if err != nil {
		return "", err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv(homeEnv)
	}

	if runtime.GOOS == consts.Darwin {
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
