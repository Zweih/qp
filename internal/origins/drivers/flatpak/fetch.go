package flatpak

import (
	"fmt"
	"qp/internal/pkgdata"
)

func fetchPackages(origin string, installDirs []string) ([]*pkgdata.PkgInfo, error) {
	pkgRefs, err := discoverPackages(installDirs)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	for _, pkgRef := range pkgRefs {
		fmt.Println(*pkgRef)
	}

	return []*pkgdata.PkgInfo{}, nil
}
