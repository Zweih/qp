package pacman

import (
	"io"
	"os"
	"strings"
	"time"
)

const bufferSize = 8192

func parseLogHistory(latestLogTime int64) (map[string]int64, int64, error) {
	path := "/var/log/pacman.log"

	file, err := os.Open(path)
	if err != nil {
		return nil, latestLogTime, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, latestLogTime, err
	}

	return readLogBackwards(file, stat.Size(), latestLogTime)
}

func readLogBackwards(
	file *os.File,
	fileSize int64,
	latestLogTime int64,
) (map[string]int64, int64, error) {
	buffer := make([]byte, bufferSize)
	installTimes := make(map[string]int64)
	removeTimes := make(map[string]int64)

	var lineBuffer strings.Builder
	pos := fileSize
	var newestTimestamp int64
	var currentTimestamp int64

	for pos > 0 {
		readSize := min(int64(bufferSize), pos)
		pos -= readSize

		numBytes, err := file.ReadAt(buffer[:readSize], pos)
		if err != nil && err != io.EOF {
			return nil, latestLogTime, err
		}

		for i := numBytes - 1; i >= 0; i-- {
			if buffer[i] == '\n' {
				if lineBuffer.Len() > 0 {
					line := reverseLine(lineBuffer.String())

					currentTimestamp = processLine(line, installTimes, removeTimes)
					if newestTimestamp == 0 {
						newestTimestamp = currentTimestamp
					}

					if currentTimestamp > 0 && currentTimestamp <= latestLogTime {
						return installTimes, newestTimestamp, nil
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

	return installTimes, newestTimestamp, nil
}

func processLine(
	line string,
	installTimes map[string]int64,
	removeTimes map[string]int64,
) int64 {
	parts := strings.SplitN(line, " ", 5)
	if len(parts) != 5 {
		return 0
	}

	rawTime := parts[0]
	t, err := time.Parse("2006-01-02T15:04:05-0700", rawTime[1:len(rawTime)-1])
	if err != nil {
		return 0
	}

	timestamp := t.Unix()

	if parts[1] != "[ALPM]" {
		return timestamp
	}

	var currentMap, oppositeMap map[string]int64
	action := parts[2]

	switch action {
	case "installed":
		currentMap, oppositeMap = installTimes, removeTimes
	case "removed":
		currentMap, oppositeMap = removeTimes, installTimes
	default:
		return timestamp
	}

	name := parts[3]

	if _, exists := currentMap[name]; exists {
		return timestamp
	}

	if _, existsInOpposite := oppositeMap[name]; existsInOpposite {
		delete(oppositeMap, name)
		return timestamp
	}

	currentMap[name] = timestamp
	return timestamp
}

func reverseLine(line string) string {
	chars := []rune(line)

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}
