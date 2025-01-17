package pkgdata

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func FetchPackages() ([]PackageInfo, error) {
	// expac queries Arch Linux + Arch-based package DBs
	cmd := exec.Command("expac", "--timefmt=%Y-%m-%d %T", "%l\t%n\t%w\t%m")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error starting command: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting command: %w", err)
	}

	var packages []PackageInfo
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		// check for correct field format, skip if not
		if len(fields) != 4 {
			continue
		}

		timestampStr, name, reason, sizeStr := fields[0], fields[1], fields[2], fields[3]
		timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			continue
		}

		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			size = 0
		}

		packages = append(packages, PackageInfo{
			Timestamp: timestamp,
			Name:      name,
			Reason:    reason,
			Size:      size,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading command output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("comand exited with error: %w", err)
	}

	return packages, nil
}
