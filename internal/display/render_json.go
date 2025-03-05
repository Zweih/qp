package display

import (
	"encoding/json"
	"fmt"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

func (o *OutputManager) renderJson(pkgs []pkgdata.PackageInfo, columnNames []string) {
	filteredPackages := make([]pkgdata.PackageInfoJson, len(pkgs))
	for i, pkg := range pkgs {
		filteredPackages[i] = getJsonValues(pkg, columnNames)
	}

	jsonOutput, err := json.MarshalIndent(filteredPackages, "", "  ")
	if err != nil {
		o.writeLine(fmt.Sprintf("Error genereating JSON output: %v", err))
	}

	o.writeLine(string(jsonOutput))
}

func getJsonValues(pkg pkgdata.PackageInfo, columnNames []string) pkgdata.PackageInfoJson {
	filteredPackage := pkgdata.PackageInfoJson{}

	for _, columnName := range columnNames {
		switch columnName {
		case consts.Date:
			filteredPackage.Timestamp = &pkg.Timestamp
		case consts.Name:
			filteredPackage.Name = pkg.Name
		case consts.Reason:
			filteredPackage.Reason = pkg.Reason
		case consts.Size:
			filteredPackage.Size = pkg.Size // return in bytes for json
		case consts.Version:
			filteredPackage.Version = pkg.Version
		case consts.Depends:
			filteredPackage.Depends = pkg.Depends
		case consts.RequiredBy:
			filteredPackage.RequiredBy = pkg.RequiredBy
		case consts.Provides:
			filteredPackage.Provides = pkg.Provides
		}
	}

	return filteredPackage
}
