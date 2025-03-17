package display

import (
	"encoding/json"
	"fmt"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

func (o *OutputManager) renderJson(pkgs []pkgdata.PackageInfo, fields []consts.FieldType) {
	filteredPackages := make([]pkgdata.PackageInfoJson, len(pkgs))
	for i, pkg := range pkgs {
		filteredPackages[i] = getJsonValues(pkg, fields)
	}

	jsonOutput, err := json.MarshalIndent(filteredPackages, "", "  ")
	if err != nil {
		o.writeLine(fmt.Sprintf("Error genereating JSON output: %v", err))
	}

	o.writeLine(string(jsonOutput))
}

func getJsonValues(pkg pkgdata.PackageInfo, fields []consts.FieldType) pkgdata.PackageInfoJson {
	filteredPackage := pkgdata.PackageInfoJson{}

	for _, field := range fields {
		switch field {
		case consts.FieldDate:
			filteredPackage.Timestamp = &pkg.Timestamp
		case consts.FieldName:
			filteredPackage.Name = pkg.Name
		case consts.FieldReason:
			filteredPackage.Reason = pkg.Reason
		case consts.FieldSize:
			filteredPackage.Size = pkg.Size // return in bytes for json
		case consts.FieldVersion:
			filteredPackage.Version = pkg.Version
		case consts.FieldDepends:
			filteredPackage.Depends = pkg.Depends
		case consts.FieldRequiredBy:
			filteredPackage.RequiredBy = pkg.RequiredBy
		case consts.FieldProvides:
			filteredPackage.Provides = pkg.Provides
		case consts.FieldConflicts:
			filteredPackage.Conflicts = pkg.Conflicts
		case consts.FieldArch:
			filteredPackage.Arch = pkg.Arch
		}
	}

	return filteredPackage
}
