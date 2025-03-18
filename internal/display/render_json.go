package display

import (
	"encoding/json"
	"fmt"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

func (o *OutputManager) renderJson(pkgs []pkgdata.PackageInfo, fields []consts.FieldType) {
	if !isAllFields(fields) {
		pkgs = selectJsonFields(pkgs, fields)
	}

	jsonOutput, err := json.MarshalIndent(pkgs, "", "  ")
	if err != nil {
		o.writeLine(fmt.Sprintf("Error genereating JSON output: %v", err))
	}

	o.writeLine(string(jsonOutput))
}

// quick check to verify if we should select fields at all
func isAllFields(fields []consts.FieldType) bool {
	if len(fields) != len(consts.ValidFields) {
		return false
	}

	for _, field := range fields {
		for _, validField := range consts.ValidFields {
			if field != validField {
				return false
			}
		}
	}

	return true
}

func selectJsonFields(
	pkgs []pkgdata.PackageInfo,
	fields []consts.FieldType,
) []pkgdata.PackageInfo {
	filteredPackages := make([]pkgdata.PackageInfo, len(pkgs))
	for i, pkg := range pkgs {
		filteredPackages[i] = getJsonValues(pkg, fields)
	}

	return filteredPackages
}

func getJsonValues(pkg pkgdata.PackageInfo, fields []consts.FieldType) pkgdata.PackageInfo {
	filteredPackage := pkgdata.PackageInfo{}

	for _, field := range fields {
		switch field {
		case consts.FieldDate:
			filteredPackage.Timestamp = pkg.Timestamp
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
