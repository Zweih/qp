package brew

import (
	"fmt"
	"os"
)

func getPrefix() (string, error) {
	possiblePrefixes := []string{
		armMacPrefix,
		x86MacDetectPrefix,
		linuxPrefix,
	}

	for _, path := range possiblePrefixes {
		if _, err := os.Stat(path); err == nil {
			return normalizeMac(path), nil
		}
	}

	return "", fmt.Errorf("no valid Homebrew prefix found")
}

func normalizeMac(detectedPrefix string) string {
	if detectedPrefix == x86MacDetectPrefix {
		return x86MacPrefix
	}

	return detectedPrefix
}
