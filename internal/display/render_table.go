package display

import (
	"bytes"
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strings"
	"text/tabwriter"
)

type tableContext struct {
	DateFormat string
}

var columnHeaders = map[consts.FieldType]string{
	consts.FieldDate:        "DATE",
	consts.FieldBuildDate:   "BUILD DATE",
	consts.FieldName:        "NAME",
	consts.FieldReason:      "REASON",
	consts.FieldSize:        "SIZE",
	consts.FieldVersion:     "VERSION",
	consts.FieldOrigin:      "ORIGIN",
	consts.FieldArch:        "ARCH",
	consts.FieldLicense:     "LICENSE",
	consts.FieldPkgBase:     "PKGBASE",
	consts.FieldDescription: "DESCRIPTION",
	consts.FieldUrl:         "URL",
	consts.FieldValidation:  "VALIDATION",
	consts.FieldPkgType:     "PKGTYPE",
	consts.FieldPackager:    "PACKAGER",
	consts.FieldGroups:      "GROUPS",
	consts.FieldConflicts:   "CONFLICTS",
	consts.FieldReplaces:    "REPLACES",
	consts.FieldDepends:     "DEPENDS",
	consts.FieldOptDepends:  "OPT DEPENDS",
	consts.FieldRequiredBy:  "REQUIRED BY",
	consts.FieldOptionalFor: "OPTIONAL FOR",
	consts.FieldProvides:    "PROVIDES",
}

// displays data in tab format
func (o *OutputManager) renderTable(
	pkgPtrs []*pkgdata.PkgInfo,
	fields []consts.FieldType,
	showFullTimestamp bool,
	hasNoHeaders bool,
) {
	o.clearProgress()

	dateFormat := consts.DateOnlyFormat
	if showFullTimestamp {
		dateFormat = consts.DateTimeFormat
	}

	ctx := tableContext{DateFormat: dateFormat}

	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 8, 2, ' ', 0)

	if !hasNoHeaders {
		renderHeaders(w, fields)
	}

	for _, pkg := range pkgPtrs {
		renderRows(w, pkg, fields, ctx)
	}

	w.Flush()
	o.write(buffer.String())
}

func renderHeaders(w *tabwriter.Writer, fields []consts.FieldType) {
	headers := make([]string, len(fields))
	for i, field := range fields {
		headers[i] = columnHeaders[field]
	}

	fmt.Fprintln(w, strings.Join(headers, "\t"))
}

func renderRows(
	w *tabwriter.Writer,
	pkg *pkgdata.PkgInfo,
	fields []consts.FieldType,
	ctx tableContext,
) {
	row := make([]string, len(fields))
	for i, field := range fields {
		value := getTableValue(pkg, field, ctx)
		if value == "" {
			value = "-"
		}

		row[i] = value
	}

	fmt.Fprintln(w, strings.Join(row, "\t"))
}

func getTableValue(pkg *pkgdata.PkgInfo, field consts.FieldType, ctx tableContext) string {
	switch field {
	case consts.FieldDate, consts.FieldBuildDate:
		return formatDate(pkg.GetInt(field), ctx)

	case consts.FieldSize:
		return formatSize(pkg.GetInt(field))

	case consts.FieldName, consts.FieldReason, consts.FieldVersion,
		consts.FieldOrigin, consts.FieldArch, consts.FieldLicense,
		consts.FieldUrl, consts.FieldDescription, consts.FieldValidation,
		consts.FieldPkgType, consts.FieldPkgBase, consts.FieldPackager:
		return pkg.GetString(field)

	case consts.FieldGroups:
		return strings.Join(pkg.GetStrArr(field), ", ")

	case consts.FieldConflicts, consts.FieldReplaces, consts.FieldDepends,
		consts.FieldOptDepends, consts.FieldRequiredBy, consts.FieldOptionalFor,
		consts.FieldProvides:
		relations := pkg.GetRelations(field)
		return formatRelations(relations)
	}

	return ""
}
