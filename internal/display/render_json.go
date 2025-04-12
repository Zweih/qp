package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type PkgInfoJson struct {
	InstallTimestamp int64    `json:"installTimestamp,omitempty"`
	BuildTimestamp   int64    `json:"buildTimestamp,omitempty"`
	Size             int64    `json:"size,omitempty"`
	PkgType          string   `json:"pkgtype,omitempty"`
	Name             string   `json:"name,omitempty"`
	Reason           string   `json:"reason,omitempty"`
	Version          string   `json:"version,omitempty"`
	Arch             string   `json:"arch,omitempty"`
	License          string   `json:"license,omitempty"`
	PkgBase          string   `json:"pkgbase,omitempty"`
	Description      string   `json:"description,omitempty"`
	Url              string   `json:"url,omitempty"`
	Validation       string   `json:"validation,omitempty"`
	Packager         string   `json:"packager,omitempty"`
	Groups           []string `json:"groups,omitempty"`
	Conflicts        []string `json:"conflicts,omitempty"`
	Replaces         []string `json:"replaces,omitempty"`
	Depends          []string `json:"depends,omitempty"`
	OptDepends       []string `json:"optDepends,omitempty"`
	RequiredBy       []string `json:"requiredBy,omitempty"`
	Provides         []string `json:"provides,omitempty"`
}

func (o *OutputManager) renderJson(pkgPtrs []*pkgdata.PkgInfo, fields []consts.FieldType) {
	uniqueFields := getUniqueFields(fields)
	filteredPkgPtrs := selectJsonFields(pkgPtrs, uniqueFields)

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false) // disable escaping of characters like `<`, `>`, perhaps this should be a user defined option
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(filteredPkgPtrs); err != nil {
		o.writeLine(fmt.Sprintf("Error generating JSON output: %v", err))
	}

	o.write(buffer.String())
}

func getUniqueFields(fields []consts.FieldType) []consts.FieldType {
	fieldSet := make(map[consts.FieldType]bool, len(fields))
	for _, field := range fields {
		fieldSet[field] = true
	}

	uniqueFields := make([]consts.FieldType, 0, len(fieldSet))
	for field := range fieldSet {
		uniqueFields = append(uniqueFields, field)
	}

	return uniqueFields
}

func selectJsonFields(
	pkgPtrs []*pkgdata.PkgInfo,
	fields []consts.FieldType,
) []*PkgInfoJson {
	filteredPkgPtrs := make([]*PkgInfoJson, len(pkgPtrs))
	for i, pkg := range pkgPtrs {
		filteredPkgPtrs[i] = getJsonValues(pkg, fields)
	}

	return filteredPkgPtrs
}

func getJsonValues(pkg *pkgdata.PkgInfo, fields []consts.FieldType) *PkgInfoJson {
	filteredPackage := PkgInfoJson{}

	for _, field := range fields {
		switch field {
		case consts.FieldDate:
			filteredPackage.InstallTimestamp = pkg.InstallTimestamp
		case consts.FieldBuildDate:
			filteredPackage.BuildTimestamp = pkg.BuildTimestamp
		case consts.FieldPkgType:
			filteredPackage.PkgType = pkgTypeToString(pkg.PkgType)
		case consts.FieldName:
			filteredPackage.Name = pkg.Name
		case consts.FieldReason:
			filteredPackage.Reason = pkg.Reason
		case consts.FieldSize:
			filteredPackage.Size = pkg.Size // return in bytes for json
		case consts.FieldVersion:
			filteredPackage.Version = pkg.Version
		case consts.FieldArch:
			filteredPackage.Arch = pkg.Arch
		case consts.FieldLicense:
			filteredPackage.License = pkg.License
		case consts.FieldPkgBase:
			filteredPackage.PkgBase = pkg.PkgBase
		case consts.FieldDescription:
			filteredPackage.Description = pkg.Description
		case consts.FieldUrl:
			filteredPackage.Url = pkg.Url
		case consts.FieldGroups:
			filteredPackage.Groups = pkg.Groups
		case consts.FieldValidation:
			filteredPackage.Validation = pkg.Validation
		case consts.FieldPackager:
			filteredPackage.Packager = pkg.Packager
		case consts.FieldConflicts:
			filteredPackage.Conflicts = flattenRelations(pkg.Conflicts)
		case consts.FieldReplaces:
			filteredPackage.Replaces = flattenRelations(pkg.Replaces)
		case consts.FieldDepends:
			filteredPackage.Depends = flattenRelations(pkg.Depends)
		case consts.FieldOptDepends:
			filteredPackage.OptDepends = flattenRelations(pkg.OptDepends)
		case consts.FieldRequiredBy:
			filteredPackage.RequiredBy = flattenRelations(pkg.RequiredBy)
		case consts.FieldProvides:
			filteredPackage.Provides = flattenRelations(pkg.Provides)
		}
	}

	return &filteredPackage
}
