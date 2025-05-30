//go:build linux

package shared

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

func GetBirthTime(path string) (int64, bool, error) {
	var stat unix.Statx_t

	err := unix.Statx(unix.AT_FDCWD, path, 0, unix.STATX_BTIME, &stat)
	if err != nil {
		return 0, false, err
	}

	if stat.Mask&unix.STATX_BTIME == 0 {
		return 0, false, fmt.Errorf("birth time not available")
	}

	birthTime := time.Unix(stat.Btime.Sec, int64(stat.Btime.Nsec))
	return birthTime.Unix(), true, nil
}
