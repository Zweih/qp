package debstyle

import (
	"bytes"
	"strings"
)

func ParseStatusFields(block []byte) map[string]string {
	fields := make(map[string]string)

	for line := range bytes.SplitSeq(block, []byte("\n")) {
		// skip continuation lines
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			continue
		}

		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(string(parts[0]))
		value := strings.TrimSpace(string(parts[1]))
		fields[key] = value
	}

	return fields
}
