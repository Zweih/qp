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
	consts.FieldName:        "NAME",
	consts.FieldReason:      "REASON",
	consts.FieldSize:        "SIZE",
	consts.FieldVersion:     "VERSION",
	consts.FieldDepends:     "DEPENDS",
	consts.FieldRequiredBy:  "REQUIRED BY",
	consts.FieldProvides:    "PROVIDES",
	consts.FieldConflicts:   "CONFLICTS",
	consts.FieldArch:        "ARCH",
	consts.FieldLicense:     "LICENSE",
	consts.FieldUrl:         "URL",
	consts.FieldDescription: "DESCRIPTION",
	consts.FieldPkgBase:     "PKGBASE",
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
	for i, fields := range fields {
		row[i] = getTableValue(pkg, fields, ctx)
	}

	fmt.Fprintln(w, strings.Join(row, "\t"))
}

func getTableValue(pkg *pkgdata.PkgInfo, field consts.FieldType, ctx tableContext) string {
	switch field {
	case consts.FieldDate:
		return formatDate(pkg, ctx)
	case consts.FieldName:
		return pkg.Name
	case consts.FieldReason:
		return pkg.Reason
	case consts.FieldSize:
		return formatSize(pkg.Size)
	case consts.FieldVersion:
		return pkg.Version
	case consts.FieldDepends:
		return formatRelations(pkg.Depends)
	case consts.FieldRequiredBy:
		return formatRelations(pkgdata.GetRelationsByDepth(pkg.RequiredBy, 1))
	case consts.FieldProvides:
		return formatRelations(pkg.Provides)
	case consts.FieldConflicts:
		return formatRelations(pkg.Conflicts)
	case consts.FieldArch:
		return pkg.Arch
	case consts.FieldLicense:
		return pkg.License
	case consts.FieldUrl:
		return pkg.Url
	case consts.FieldDescription:
		return pkg.Description
	case consts.FieldPkgBase:
		return pkg.PkgBase
	default:
		return ""
	}
}
