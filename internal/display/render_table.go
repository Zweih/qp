package display

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

type tableContext struct {
	DateFormat string
}

var columnHeaders = map[consts.FieldType]string{
	consts.FieldDate:       "DATE",
	consts.FieldName:       "NAME",
	consts.FieldReason:     "REASON",
	consts.FieldSize:       "SIZE",
	consts.FieldVersion:    "VERSION",
	consts.FieldDepends:    "DEPENDS",
	consts.FieldRequiredBy: "REQUIRED BY",
	consts.FieldProvides:   "PROVIDES",
	consts.FieldConflicts:  "CONFLICTS",
	consts.FieldArch:       "ARCH",
	consts.FieldLicense:    "LICENSE",
}

// displays data in tab format
func (o *OutputManager) renderTable(
	pkgs []pkgdata.PackageInfo,
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

	for _, pkg := range pkgs {
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
	pkg pkgdata.PackageInfo,
	fields []consts.FieldType,
	ctx tableContext,
) {
	row := make([]string, len(fields))
	for i, fields := range fields {
		row[i] = getTableValue(pkg, fields, ctx)
	}

	fmt.Fprintln(w, strings.Join(row, "\t"))
}

func getTableValue(pkg pkgdata.PackageInfo, field consts.FieldType, ctx tableContext) string {
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
		return formatPackageList(pkg.Depends)
	case consts.FieldRequiredBy:
		return formatPackageList(pkg.RequiredBy)
	case consts.FieldProvides:
		return formatPackageList(pkg.Provides)
	case consts.FieldConflicts:
		return formatPackageList(pkg.Conflicts)
	case consts.FieldArch:
		return pkg.Arch
	case consts.FieldLicense:
		return pkg.License
	default:
		return ""
	}
}

// use time as parameter
func formatDate(pkg pkgdata.PackageInfo, ctx tableContext) string {
	// return pkg.Timestamp.Format(ctx.DateFormat)
	timestamp := time.Unix(pkg.Timestamp, 0)
	return timestamp.Format(ctx.DateFormat)
}

func formatPackageList(packages []string) string {
	if len(packages) == 0 {
		return "-"
	}
	return strings.Join(packages, ", ")
}

func formatSize(size int64) string {
	switch {
	case size >= consts.GB:
		return fmt.Sprintf("%.2f GB", float64(size)/(consts.GB))
	case size >= consts.MB:
		return fmt.Sprintf("%.2f MB", float64(size)/(consts.MB))
	case size >= consts.KB:
		return fmt.Sprintf("%.2f KB", float64(size)/(consts.KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
