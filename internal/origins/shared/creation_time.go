package shared

import "os"

func GetCreationTime(path string) (int64, bool, error) {
	if birthTime, reliable, err := getBirthTime(path); err == nil {
		return birthTime, reliable, nil
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, false, err
	}

	return fileInfo.ModTime().Unix(), false, nil
}
