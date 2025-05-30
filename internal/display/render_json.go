package display

import (
	"bytes"
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"

	json "github.com/goccy/go-json"
)

type PkgInfoJSON struct {
	InstalledTimestamp int64    `json:"installedTimestamp,omitempty"`
	UpdateTimestamp    int64    `json:"updateTimestamp,omitempty"`
	BuildTimestamp     int64    `json:"buildTimestamp,omitempty"`
	Size               int64    `json:"size,omitempty"`
	Name               string   `json:"name,omitempty"`
	Reason             string   `json:"reason,omitempty"`
	Version            string   `json:"version,omitempty"`
	Origin             string   `json:"origin,omitempty"`
	Arch               string   `json:"arch,omitempty"`
	License            string   `json:"license,omitempty"`
	Description        string   `json:"description,omitempty"`
	Url                string   `json:"url,omitempty"`
	Validation         string   `json:"validation,omitempty"`
	PkgType            string   `json:"pkgtype,omitempty"`
	PkgBase            string   `json:"pkgbase,omitempty"`
	Packager           string   `json:"packager,omitempty"`
	Groups             []string `json:"groups,omitempty"`
	Conflicts          []string `json:"conflicts,omitempty"`
	Replaces           []string `json:"replaces,omitempty"`
	Depends            []string `json:"depends,omitempty"`
	OptDepends         []string `json:"optDepends,omitempty"`
	RequiredBy         []string `json:"requiredBy,omitempty"`
	OptionalFor        []string `json:"optionalFor,omitempty"`
	Provides           []string `json:"provides,omitempty"`
}

func (o *OutputManager) renderJSON(pkgPtrs []*pkgdata.PkgInfo, fields []consts.FieldType) {
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
) []*PkgInfoJSON {
	filteredPkgPtrs := make([]*PkgInfoJSON, len(pkgPtrs))
	for i, pkg := range pkgPtrs {
		filteredPkgPtrs[i] = getJsonValues(pkg, fields)
	}

	return filteredPkgPtrs
}

func getJsonValues(pkg *pkgdata.PkgInfo, fields []consts.FieldType) *PkgInfoJSON {
	filteredPackage := PkgInfoJSON{}

	for _, field := range fields {
		switch field {
		case consts.FieldInstalled:
			filteredPackage.InstalledTimestamp = pkg.GetInt(field)
		case consts.FieldUpdated:
			filteredPackage.UpdateTimestamp = pkg.GetInt(field)
		case consts.FieldBuilt:
			filteredPackage.BuildTimestamp = pkg.GetInt(field)
		case consts.FieldSize:
			filteredPackage.Size = pkg.GetInt(field)
		case consts.FieldName:
			filteredPackage.Name = pkg.GetString(field)
		case consts.FieldReason:
			filteredPackage.Reason = pkg.GetString(field)
		case consts.FieldVersion:
			filteredPackage.Version = pkg.GetString(field)
		case consts.FieldOrigin:
			filteredPackage.Origin = pkg.GetString(field)
		case consts.FieldArch:
			filteredPackage.Arch = pkg.GetString(field)
		case consts.FieldLicense:
			filteredPackage.License = pkg.GetString(field)
		case consts.FieldDescription:
			filteredPackage.Description = pkg.GetString(field)
		case consts.FieldUrl:
			filteredPackage.Url = pkg.GetString(field)
		case consts.FieldValidation:
			filteredPackage.Validation = pkg.GetString(field)
		case consts.FieldPkgType:
			filteredPackage.PkgType = pkg.GetString(field)
		case consts.FieldPkgBase:
			filteredPackage.PkgBase = pkg.GetString(field)
		case consts.FieldPackager:
			filteredPackage.Packager = pkg.GetString(field)
		case consts.FieldGroups:
			filteredPackage.Groups = pkg.Groups
		case consts.FieldConflicts:
			filteredPackage.Conflicts = flattenRelations(pkg.GetRelations(field))
		case consts.FieldReplaces:
			filteredPackage.Replaces = flattenRelations(pkg.GetRelations(field))
		case consts.FieldDepends:
			filteredPackage.Depends = flattenRelations(pkg.GetRelations(field))
		case consts.FieldOptDepends:
			filteredPackage.OptDepends = flattenRelations(pkg.GetRelations(field))
		case consts.FieldRequiredBy:
			filteredPackage.RequiredBy = flattenRelations(pkg.GetRelations(field))
		case consts.FieldOptionalFor:
			filteredPackage.OptionalFor = flattenRelations(pkg.GetRelations(field))
		case consts.FieldProvides:
			filteredPackage.Provides = flattenRelations(pkg.GetRelations(field))
		}
	}

	return &filteredPackage
}
