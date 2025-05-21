package brew

import (
	"fmt"
	"os"
	"path/filepath"
)

func getModTime(prefix string, subPath string) (int64, error) {
	path := filepath.Join(prefix, subPath)
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read Cellar: %w", err)
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

	cellarInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to stat Cellar: %w", err)
	}

	modTime := cellarInfo.ModTime().Unix()
	if modTime > latest {
		latest = modTime
	}

	return latest, nil
}
