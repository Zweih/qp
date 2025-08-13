package apk

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func loadInstallReasons() (map[string]bool, error) {
	file, err := os.Open(apkWorldPath)
	if err != nil {
		return map[string]bool{}, fmt.Errorf("failed to open world file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return map[string]bool{}, fmt.Errorf("failed to read world file: %w", err)
	}

	reasonMap := make(map[string]bool)

	for line := range bytes.SplitSeq(data, []byte("\n")) {
		if len(line) < 1 {
			continue
		}

		reasonMap[parseWorldEntry(string(line))] = true
	}

	return reasonMap, nil
}

// TODO: perhaps we can iterate byte by byte like we do for pacman
func parseWorldEntry(entry string) string {
	// remove version constraint
	for _, op := range []string{">=", "<=", "!=", ">", "<", "=", "~"} {
		if idx := strings.Index(entry, op); idx != -1 {
			entry = entry[:idx]
			break
		}
	}

	// remove repo tag
	if idx := strings.Index(entry, "@"); idx != -1 {
		entry = entry[:idx]
	}

	return strings.TrimSpace(entry)
}
