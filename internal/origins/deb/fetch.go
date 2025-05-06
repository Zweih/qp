package deb

import (
	"fmt"
	"io"
	"os"
	"qp/internal/pkgdata"
)

func fetchPackages(origin string, reasonMap map[string]string) ([]*pkgdata.PkgInfo, error) {
	file, err := os.Open(dpkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open status file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read status file: %w", err)
	}

	return parseStatusFile(data, origin, reasonMap)
}
