package flatpak

import (
	"bytes"
	"os"
	"strings"
)

const sectionDesktopEntry = "[Desktop Entry]"

func parseDesktopFile(pkgRef *PkgRef) {
	if pkgRef.DesktopPath == "" {
		return
	}

	data, err := os.ReadFile(pkgRef.DesktopPath)
	if err != nil {
	}

	var currentSection string
	for rawLine := range bytes.SplitSeq(data, []byte("\n")) {
		if len(rawLine) > 0 {
			line := strings.TrimSpace(string(rawLine))
			if line == sectionDesktopEntry {
				currentSection = line
				continue
			}

			if currentSection == sectionDesktopEntry {
				parts := strings.Split(line, "=")
				if len(parts) != 2 {
					continue
				}

				key, value := parts[0], parts[1]
				if key == "Name" {
					pkgRef.Pkg.Title = strings.TrimSpace(value)
				}
			}
		}
	}
}
