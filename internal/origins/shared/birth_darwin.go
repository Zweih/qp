//go:build darwin

package shared

import (
	"os"
	"syscall"
	"time"
)

func getBirthTime(path string) (int64, bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, false, err
	}

	stat := fileInfo.Sys().(*syscall.Stat_t)
	birthTime := time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
	return birthTime.Unix(), true, nil
}
