package display

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

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

	ctx := displayContext{DateFormat: dateFormat}

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

func renderRows(w *tabwriter.Writer, pkg pkgdata.PackageInfo, columnNames []string, ctx displayContext) {
	row := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		row[i] = GetColumnTableValue(pkg, columnName, ctx)
	}

	fmt.Fprintln(w, strings.Join(row, "\t"))
}
