package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

func (o *OutputManager) renderJson(pkgPtrs []*pkgdata.PkgInfo, fields []consts.FieldType) {
	if isAllFields, uniqueFields := getUniqueFields(fields); isAllFields {
		pkgPtrs = selectJsonFields(pkgPtrs, uniqueFields)
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false) // disable escaping of characters like `<`, `>`, perhaps this should be a user defined option
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(pkgPtrs); err != nil {
		o.writeLine(fmt.Sprintf("Error genereating JSON output: %v", err))
	}

	o.writeLine(buffer.String())
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
	pkgPtrs []*pkgdata.PkgInfo,
	fields []consts.FieldType,
) []*pkgdata.PkgInfo {
	filteredPkgPtrs := make([]*pkgdata.PkgInfo, len(pkgPtrs))
	for i, pkg := range pkgPtrs {
		filteredPkgPtrs[i] = getJsonValues(pkg, fields)
	}

	return filteredPkgPtrs
}

func getJsonValues(pkg *pkgdata.PkgInfo, fields []consts.FieldType) *pkgdata.PkgInfo {
	filteredPackage := pkgdata.PkgInfo{}

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
		case consts.FieldUrl:
			filteredPackage.Url = pkg.Url
		}
	}

	return &filteredPackage
}
