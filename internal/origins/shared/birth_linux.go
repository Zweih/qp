//go:build linux

package shared

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

const (
	EXT4_SUPER_MAGIC      = 0xEF53
	BTRFS_SUPER_MAGIC     = 0x9123683E
	XFS_SUPER_MAGIC       = 0x58465342
	ZFS_SUPER_MAGIC       = 0x2FC12FC1
	TMPFS_MAGIC           = 0x01021994
	OVERLAYFS_SUPER_MAGIC = 0x794C7630
	FUSE_SUPER_MAGIC      = 0x65735546
)

func getFilesystemType(path string) (int64, error) {
	var statfs unix.Statfs_t

	err := unix.Statfs(path, &statfs)
	if err != nil {
		return 0, err
	}

	return statfs.Type, nil
}

func fsSupportsBtime(path string) bool {
	fsType, err := getFilesystemType(path)
	if err != nil {
		return false
	}

	switch fsType {
	case EXT4_SUPER_MAGIC, BTRFS_SUPER_MAGIC, XFS_SUPER_MAGIC, ZFS_SUPER_MAGIC:
		return true
	case TMPFS_MAGIC, OVERLAYFS_SUPER_MAGIC, FUSE_SUPER_MAGIC:
		return false
	default:
		return false
	}
}

func getBirthTime(path string) (int64, bool, error) {
	if !fsSupportsBtime(path) {
		return 0, false, fmt.Errorf("filesystem doesn't support birth time")
	}

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
