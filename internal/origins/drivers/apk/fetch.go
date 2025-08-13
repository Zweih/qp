package apk

import (
	"fmt"
	"io"
	"os"
	"qp/internal/pkgdata"
)

func fetchPackages(origin string) ([]*pkgdata.PkgInfo, error) {
	file, err := os.Open(apkDbPath)
	if err != nil {
		return []*pkgdata.PkgInfo{}, fmt.Errorf("failed to open apk database: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return []*pkgdata.PkgInfo{}, fmt.Errorf("failed to read apk database: %w", err)
	}

	reasonMap, err := loadInstallReasons()
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	fmt.Println(reasonMap)

	return parseInstalledFile(data, origin, reasonMap)
}
