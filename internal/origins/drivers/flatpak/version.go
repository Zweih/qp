package flatpak

import (
	"bytes"
	"os"
)

func extractVersion(deployPath string) (string, error) {
	data, err := os.ReadFile(deployPath)
	if err != nil {
		return "", err
	}

	pattern := []byte(appdataVersion)
	index := bytes.Index(data, pattern)
	if index == -1 {
		return "", nil
	}

	start := index + len(pattern)
	if start >= len(data) {
		return "", nil
	}

	if start < len(data) && data[start] == 0 {
		start++
	}

	var version []byte
	for i := start; i < len(data); i++ {
		b := data[i]
		if b == 0 || b < 32 || b > 126 {
			break
		}

		version = append(version, b)
	}

	return string(version), nil
}
