package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type PkgInfoJson struct {
	Timestamp   int64    `json:"timestamp,omitempty"`
	Size        int64    `json:"size,omitempty"`
	Name        string   `json:"name,omitempty"`
	Reason      string   `json:"reason,omitempty"`
	Version     string   `json:"version,omitempty"`
	Arch        string   `json:"arch,omitempty"`
	License     string   `json:"license,omitempty"`
	Url         string   `json:"url,omitempty"`
	Description string   `json:"description,omitempty"`
	PkgBase     string   `json:"pkgbase,omitempty"`
	Depends     []string `json:"depends,omitempty"`
	RequiredBy  []string `json:"requiredBy,omitempty"`
	Provides    []string `json:"provides,omitempty"`
	Conflicts   []string `json:"conflicts,omitempty"`
}

func (o *OutputManager) renderJson(pkgPtrs []*pkgdata.PkgInfo, fields []consts.FieldType) {
	uniqueFields := getUniqueFields(fields)
	filteredPkgPtrs := selectJsonFields(pkgPtrs, uniqueFields)

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false) // disable escaping of characters like `<`, `>`, perhaps this should be a user defined option
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(filteredPkgPtrs); err != nil {
		o.writeLine(fmt.Sprintf("Error genereating JSON output: %v", err))
	}

	o.writeLine(buffer.String())
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
			filteredPackage.Timestamp = pkg.Timestamp
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
		case consts.FieldUrl:
			filteredPackage.Url = pkg.Url
		case consts.FieldDescription:
			filteredPackage.Description = pkg.Description
		case consts.FieldDepends:
			filteredPackage.Depends = flattenRelations(pkg.Depends)
		case consts.FieldRequiredBy:
			filteredPackage.RequiredBy = flattenRelations(pkg.RequiredBy)
		case consts.FieldProvides:
			filteredPackage.Provides = flattenRelations(pkg.Provides)
		case consts.FieldConflicts:
			filteredPackage.Conflicts = flattenRelations(pkg.Conflicts)
		}
	}

	return &filteredPackage
}

func flattenRelations(relations []pkgdata.Relation) []string {
	relationOutputs := make([]string, 0, len(relations))

	for _, rel := range relations {
		if rel.Operator == pkgdata.OpNone {
			relationOutputs = append(relationOutputs, rel.Name)
		} else {
			op := relationOpToString(rel.Operator)
			relationOutputs = append(relationOutputs, fmt.Sprintf("%s%s%s", rel.Name, op, rel.Version))
		}
	}

	return relationOutputs
}

func relationOpToString(op pkgdata.RelationOp) string {
	switch op {
	case pkgdata.OpEqual:
		return "="
	case pkgdata.OpLess:
		return "<"
	case pkgdata.OpLessEqual:
		return "<="
	case pkgdata.OpGreater:
		return ">"
	case pkgdata.OpGreaterEqual:
		return ">="
	default:
		return ""
	}
}
