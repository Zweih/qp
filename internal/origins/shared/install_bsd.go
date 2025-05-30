//go:build freebsd

package shared

import (
	"os"
	"syscall"
	"time"
)

func GetBirthTime(path string) (int64, bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, false, err
	}

	stat := fileInfo.Sys().(*syscall.Stat_t)
	birthTime := time.Unix(int64(stat.Birthtimespec.Sec), int64(stat.Birthtimespec.Nsec))
	return birthTime.Unix(), true, nil
}
