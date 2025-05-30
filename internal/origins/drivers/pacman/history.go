package pacman

import (
	"io"
	"os"
	"strings"
	"time"
)

const bufferSize = 8192

func parseLogHistory(numPkgs int) (map[string]int64, error) {
	path := "/var/log/pacman.log"

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return readLogBackwards(file, stat.Size(), numPkgs)
}

func readLogBackwards(file *os.File, fileSize int64, numPkgs int) (map[string]int64, error) {
	buffer := make([]byte, bufferSize)
	installTimes := make(map[string]int64)
	removeTimes := make(map[string]int64)

	var lineBuffer strings.Builder
	pos := fileSize

	for pos > 0 {
		readSize := min(int64(bufferSize), pos)
		pos -= readSize

		numBytes, err := file.ReadAt(buffer[:readSize], pos)
		if err != nil && err != io.EOF {
			return nil, err
		}

		for i := numBytes - 1; i >= 0; i-- {
			if buffer[i] == '\n' {
				if lineBuffer.Len() > 0 {
					line := reverseLine(lineBuffer.String())

					if processLine(line, installTimes, removeTimes) && len(installTimes) >= numPkgs {
						return installTimes, nil
					}

					lineBuffer.Reset()
				}

				continue
			}

			lineBuffer.WriteByte(buffer[i])
		}
	}

	// final line
	if lineBuffer.Len() > 0 {
		line := reverseLine(lineBuffer.String())
		processLine(line, installTimes, removeTimes)
	}

	return installTimes, nil
}

func processLine(line string, installTimes map[string]int64, removeTimes map[string]int64) bool {
	parts := strings.SplitN(line, " ", 5)
	if len(parts) != 5 {
		return false
	}

	if parts[1] != "[ALPM]" {
		return false
	}

	var currentMap, oppositeMap map[string]int64
	action := parts[2]
	var isInstalled bool

	switch action {
	case "installed":
		currentMap, oppositeMap = installTimes, removeTimes
		isInstalled = true
	case "removed":
		currentMap, oppositeMap = removeTimes, installTimes
	default:
		return false
	}

	name := parts[3]

	if _, exists := currentMap[name]; exists {
		return false
	}

	if _, existsInOpposite := oppositeMap[name]; existsInOpposite {
		delete(oppositeMap, name)
		return false
	}

	rawTime := parts[0]
	timestamp, err := time.Parse("2006-01-02T15:04:05-0700", rawTime[1:len(rawTime)-1])
	if err != nil {
		return false
	}

	currentMap[name] = timestamp.Unix()
	return isInstalled
}

func reverseLine(line string) string {
	chars := []rune(line)

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}
