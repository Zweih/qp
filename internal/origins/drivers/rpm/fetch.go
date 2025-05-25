package rpm

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"

	rpmdb "github.com/Zweih/go-rpmdb/pkg"
	_ "github.com/glebarez/sqlite"
)

func fetchPackages(origin string, path string) ([]*pkgdata.PkgInfo, error) {
	reasonMap, err := loadInstallReasons()
	if err != nil {
		fmt.Printf("WARNING: rpm install reasons will display as 'explicit': %v \n", err)
	}

	db, err := rpmdb.Open(path)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	rpmPkgList, err := db.ListPackages()
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}
	defer db.Close()

	var pkgs []*pkgdata.PkgInfo

	for _, rpmPkg := range rpmPkgList {
		group := rpmPkg.Group
		if group == groupUnspecified {
			group = ""
		}

		reason := consts.ReasonExplicit
		if historyReason, exists := reasonMap[rpmPkg.Name]; exists {
			reason = historyReason
		}

		pkg := &pkgdata.PkgInfo{
			InstallTimestamp: int64(rpmPkg.InstallTime),
			BuildTimestamp:   int64(rpmPkg.BuildTime),
			Name:             rpmPkg.Name,
			Version:          parseVersion(rpmPkg.Epoch, rpmPkg.Version),
			Arch:             rpmPkg.Arch,
			Size:             int64(rpmPkg.Size),
			License:          rpmPkg.License,
			Origin:           origin,
			Description:      rpmPkg.Summary,
			Packager:         rpmPkg.Packager,
			Url:              rpmPkg.URL,
			Groups:           []string{group},
			Reason:           reason,
			Conflicts:        parseRelationList(rpmPkg.Conflicts),
			Replaces:         parseRelationList(rpmPkg.Obsoletes),
			Depends:          parseRelationList(rpmPkg.Requires),
			Provides:         parseRelationList(rpmPkg.Provides),
		}

		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}
