package rpm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	out "qp/internal/display"
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
			if candidate == sqliteDbFile {
				if err := checkSqlite(); err != nil {
					out.WriteLine(fmt.Sprintf("WARNING: %v", err))
					continue
				}
			}

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

func checkSqlite() error {
	_, err := exec.LookPath("sqlite3")
	if err != nil {
		return fmt.Errorf("RPM databases detected but sqlite3 command not found - please install sqlite to view RPM origin")
	}

	return nil
}
