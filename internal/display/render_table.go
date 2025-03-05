package display

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

type tableContext struct {
	DateFormat string
}

var columnHeaders = map[string]string{
	consts.Date:       "DATE",
	consts.Name:       "NAME",
	consts.Reason:     "REASON",
	consts.Size:       "SIZE",
	consts.Version:    "VERSION",
	consts.Depends:    "DEPENDS",
	consts.RequiredBy: "REQUIRED BY",
	consts.Provides:   "PROVIDES",
}

// displays data in tab format
func (o *OutputManager) renderTable(
	pkgs []pkgdata.PackageInfo,
	columnNames []string,
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
		renderHeaders(w, columnNames)
	}

	for _, pkg := range pkgs {
		renderRows(w, pkg, columnNames, ctx)
	}

	w.Flush()
	o.write(buffer.String())
}

func renderHeaders(w *tabwriter.Writer, columnNames []string) {
	headers := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		headers[i] = columnHeaders[columnName]
	}

	fmt.Fprintln(w, strings.Join(headers, "\t"))
}

func renderRows(w *tabwriter.Writer, pkg pkgdata.PackageInfo, columnNames []string, ctx tableContext) {
	row := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		row[i] = getTableValue(pkg, columnName, ctx)
	}

	fmt.Fprintln(w, strings.Join(row, "\t"))
}

func getTableValue(pkg pkgdata.PackageInfo, columnName string, ctx tableContext) string {
	switch columnName {
	case consts.Date:
		return formatDate(pkg, ctx)
	case consts.Name:
		return pkg.Name
	case consts.Reason:
		return pkg.Reason
	case consts.Size:
		return formatSize(pkg.Size)
	case consts.Version:
		return pkg.Version
	case consts.Depends:
		return formatPackageList(pkg.Depends)
	case consts.RequiredBy:
		return formatPackageList(pkg.RequiredBy)
	case consts.Provides:
		return formatPackageList(pkg.Provides)
	default:
		return ""
	}
}

// use time as parameter
func formatDate(pkg pkgdata.PackageInfo, ctx tableContext) string {
	return pkg.Timestamp.Format(ctx.DateFormat)
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
