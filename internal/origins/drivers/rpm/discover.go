package rpm

import (
	"os"
	"path/filepath"
)

func detectRpmDatabase(basePath string) string {
	if _, err := os.Stat(basePath); err != nil {
		return ""
	}

	candidates := []string{
		sqliteDbFile,
		ndbDbFile,
		berkeleyDbFile,
	}

	for _, candidate := range candidates {
		fullPath := filepath.Join(basePath, candidate)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return ""
}

func findHistoryDb() string {
	candidates := []string{
		dnfHistoryPath,
		yumHistoryPath,
	}

	for _, pattern := range candidates {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return matches[len(matches)-1]
		}
	}

	matches, err := filepath.Glob(yumHistoryPattern)
	if err == nil && len(matches) > 0 {
		return matches[len(matches)-1] // most recent (lexographically)
	}

	return ""
}
