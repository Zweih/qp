package deb

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
)

func getInstallTime(pkg *pkgdata.PkgInfo) error {
	installDate, buildDate, ok := getModTime(pkg.Name)
	if ok {
		pkg.UpdateTimestamp = installDate
		pkg.BuildTimestamp = buildDate
		return nil
	}

	suffix := ":" + pkg.Arch

	installDate, buildDate, ok = getModTime(pkg.Name + suffix)
	if ok {
		pkg.UpdateTimestamp = installDate
		pkg.BuildTimestamp = buildDate
		return nil
	}

	return fmt.Errorf("Could not find .list file for %s", pkg.Name)
}

func getModTime(fileName string) (int64, int64, bool) {
	dirInfo, err := os.Stat(filepath.Join(pkgModRoot, fileName+listExt))
	if err == nil {
		installDate := dirInfo.ModTime().Unix()
		dirInfo, _ = os.Stat(filepath.Join(pkgModRoot, fileName+md5SumsExt))
		buildDate := dirInfo.ModTime().Unix()
		return installDate, buildDate, true
	}

	return 0, 0, false
}
