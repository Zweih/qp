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
	Env                string   `json:"env,omitempty"`
	PkgType            string   `json:"pkgtype,omitempty"`
	PkgBase            string   `json:"pkgbase,omitempty"`
	Packager           string   `json:"packager,omitempty"`
	Groups             []string `json:"groups,omitempty"`
	AlsoIn             []string `json:"alsoIn,omitempty"`
	OtherEnvs          []string `json:"otherEnvs,omitempty"`
	Conflicts          []string `json:"conflicts,omitempty"`
	Replaces           []string `json:"replaces,omitempty"`
	Depends            []string `json:"depends,omitempty"`
	OptDepends         []string `json:"optDepends,omitempty"`
	RequiredBy         []string `json:"requiredBy,omitempty"`
	OptionalFor        []string `json:"optionalFor,omitempty"`
	Provides           []string `json:"provides,omitempty"`
}

func (o *OutputManager) renderJSON(pkgs []*pkgdata.PkgInfo, fields []consts.FieldType) {
	uniqueFields := getUniqueFields(fields)
	filteredPkgs := selectJsonFields(pkgs, uniqueFields)

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false) // disable escaping of characters like `<`, `>`, perhaps this should be a user defined option
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(filteredPkgs); err != nil {
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
	pkgs []*pkgdata.PkgInfo,
	fields []consts.FieldType,
) []*PkgInfoJSON {
	filteredPkgs := make([]*PkgInfoJSON, len(pkgs))
	for i, pkg := range pkgs {
		filteredPkgs[i] = getJsonValues(pkg, fields)
	}

	return filteredPkgs
}

func getJsonValues(pkg *pkgdata.PkgInfo, fields []consts.FieldType) *PkgInfoJSON {
	filteredPkg := PkgInfoJSON{}

	for _, field := range fields {
		switch field {
		case consts.FieldInstalled:
			filteredPkg.InstalledTimestamp = pkg.GetInt(field)
		case consts.FieldUpdated:
			filteredPkg.UpdateTimestamp = pkg.GetInt(field)
		case consts.FieldBuilt:
			filteredPkg.BuildTimestamp = pkg.GetInt(field)
		case consts.FieldSize:
			filteredPkg.Size = pkg.GetInt(field)
		case consts.FieldName:
			filteredPkg.Name = pkg.GetString(field)
		case consts.FieldReason:
			filteredPkg.Reason = pkg.GetString(field)
		case consts.FieldVersion:
			filteredPkg.Version = pkg.GetString(field)
		case consts.FieldOrigin:
			filteredPkg.Origin = pkg.GetString(field)
		case consts.FieldArch:
			filteredPkg.Arch = pkg.GetString(field)
		case consts.FieldLicense:
			filteredPkg.License = pkg.GetString(field)
		case consts.FieldDescription:
			filteredPkg.Description = pkg.GetString(field)
		case consts.FieldUrl:
			filteredPkg.Url = pkg.GetString(field)
		case consts.FieldValidation:
			filteredPkg.Validation = pkg.GetString(field)
		case consts.FieldEnv:
			filteredPkg.Env = pkg.GetString(field)
		case consts.FieldPkgType:
			filteredPkg.PkgType = pkg.GetString(field)
		case consts.FieldPkgBase:
			filteredPkg.PkgBase = pkg.GetString(field)
		case consts.FieldPackager:
			filteredPkg.Packager = pkg.GetString(field)
		case consts.FieldGroups:
			filteredPkg.Groups = pkg.GetStrArr(field)
		case consts.FieldAlsoIn:
			filteredPkg.AlsoIn = pkg.GetStrArr(field)
		case consts.FieldOtherEnvs:
			filteredPkg.OtherEnvs = pkg.GetStrArr(field)
		case consts.FieldConflicts:
			filteredPkg.Conflicts = flattenRelations(pkg.GetRelations(field))
		case consts.FieldReplaces:
			filteredPkg.Replaces = flattenRelations(pkg.GetRelations(field))
		case consts.FieldDepends:
			filteredPkg.Depends = flattenRelations(pkg.GetRelations(field))
		case consts.FieldOptDepends:
			filteredPkg.OptDepends = flattenRelations(pkg.GetRelations(field))
		case consts.FieldRequiredBy:
			filteredPkg.RequiredBy = flattenRelations(pkg.GetRelations(field))
		case consts.FieldOptionalFor:
			filteredPkg.OptionalFor = flattenRelations(pkg.GetRelations(field))
		case consts.FieldProvides:
			filteredPkg.Provides = flattenRelations(pkg.GetRelations(field))
		}
	}

	return &filteredPkg
}
