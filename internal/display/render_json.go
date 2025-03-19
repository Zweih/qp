package display

import (
	"encoding/json"
	"fmt"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

func (o *OutputManager) renderJson(pkgs []pkgdata.PackageInfo, fields []consts.FieldType) {
	if isAllFields, uniqueFields := getUniqueFields(fields); isAllFields {
		pkgs = selectJsonFields(pkgs, uniqueFields)
	}

	jsonOutput, err := json.MarshalIndent(pkgs, "", "  ")
	if err != nil {
		o.writeLine(fmt.Sprintf("Error genereating JSON output: %v", err))
	}

	o.writeLine(string(jsonOutput))
}

// quick check to verify if we should select fields at all
func getUniqueFields(fields []consts.FieldType) (bool, []consts.FieldType) {
	fieldSet := make(map[consts.FieldType]bool, len(fields))
	for _, field := range fields {
		fieldSet[field] = true
	}

	uniqueFields := make([]consts.FieldType, 0, len(fieldSet))
	for field := range fieldSet {
		uniqueFields = append(uniqueFields, field)
	}

	return len(fieldSet) != len(consts.ValidFields), uniqueFields
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
		case consts.FieldLicense:
			filteredPackage.License = pkg.License
		}
	}

	return filteredPackage
}
