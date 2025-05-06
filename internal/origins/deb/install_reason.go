package deb

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"qp/internal/consts"
)

// for dpkg only systems and minimal ubuntu docker containers, this file does not exist or is empty. we use a fallback after dependency resolution in those cases to infer the reason
func loadInstallReasons() (map[string]string, error) {
	file, err := os.Open(installReasonPath)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to read extended_states: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to load extended_states: %w", err)
	}

	if len(data) == 0 {
		return map[string]string{}, fmt.Errorf("extended_states file is empty")
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
			reasonMap[name] = consts.ReasonDependency
		}
	}

	return reasonMap, nil
}
