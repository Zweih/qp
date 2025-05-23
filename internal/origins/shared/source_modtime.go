package shared

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetModTime(path string) (int64, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var latest int64
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := os.Stat(filepath.Join(path, entry.Name()))
		if err != nil {
			continue
		}

		modTime := info.ModTime().Unix()
		if modTime > latest {
			latest = modTime
		}
	}

	parentInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to stat %s: %w", path, err)
	}

	modTime := parentInfo.ModTime().Unix()
	if modTime > latest {
		latest = modTime
	}

	return latest, nil
}
