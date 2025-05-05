package deb

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// only works for APT systems, we'll need to add more heuristics to get reasons for dpkg only systems
func loadInstallReasons() (map[string]string, error) {
	file, err := os.Open(installReasonPath)
	if err != nil {
		return map[string]string{}, nil
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read install reason file: %w", err)
	}

	reasonMap := make(map[string]string)

	for block := range bytes.SplitSeq(data, []byte("\n\n")) {
		var name string
		var isAuto bool

		for line := range bytes.SplitSeq(block, []byte("\n")) {
			if bytes.HasPrefix(line, []byte(packagePrefix)) {
				name = string(bytes.TrimSpace(line[len(packagePrefix):]))
			} else if bytes.HasPrefix(line, []byte(autoInstallPrefix)) {
				val := string(bytes.TrimSpace(line[len(autoInstallPrefix):]))
				isAuto = val == "1"
			}
		}

		if name != "" && isAuto {
			reasonMap[name] = "dependency"
		}
	}

	return reasonMap, nil
}
