package shared

import (
	"io/fs"
	"path/filepath"
)

func GetInstallSize(dir string) (int64, error) {
	var total int64

	err := filepath.WalkDir(dir, func(_ string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !dir.IsDir() {
			info, err := dir.Info()
			if err != nil {
				return err
			}

			total += info.Size()
		}

		return nil
	})

	return total, err
}
