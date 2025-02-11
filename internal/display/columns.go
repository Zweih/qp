package display

import (
	"fmt"
	"strings"
	"yaylog/internal/pkgdata"
)

type DisplayContext struct {
	DateFormat string
}

type PackageInfo = pkgdata.PackageInfo

type Column struct {
	Name      string
	Header    string
	Getter    func(pkg PackageInfo, ctx DisplayContext) string
	IsDefault bool
}

var allColumns = []Column{
	newDateColumn(),
	newSimpleColumn("name", "NAME", func(pkg PackageInfo) string { return pkg.Name }, true),
	newSimpleColumn("version", "VERSION", func(pkg PackageInfo) string { return pkg.Version }, false),
	newSimpleColumn("reason", "REASON", func(pkg PackageInfo) string { return pkg.Reason }, true),
	newSimpleColumn("size", "SIZE", func(pkg PackageInfo) string { return formatSize(pkg.Size) }, true),
}

func newDateColumn() Column {
	return Column{
		Name:      "date",
		Header:    "DATE",
		IsDefault: true,
		Getter: func(pkg PackageInfo, ctx DisplayContext) string {
			return pkg.Timestamp.Format(ctx.DateFormat)
		},
	}
}

func newSimpleColumn(
	name string,
	header string,
	getter func(pkg PackageInfo) string,
	isDefault bool,
) Column {
	return Column{
		Name:   name,
		Header: header,
		Getter: func(pkg PackageInfo, _ DisplayContext) string {
			return getter(pkg)
		},
		IsDefault: isDefault,
	}
}

func GetActiveColumns(withDefaultColumns bool, optionalColumns []string) ([]Column, error) {
	var activeColumns []Column

	if withDefaultColumns {
		activeColumns = filterDefaultColumns()
	}

	unknownColumns := []string{}

	for _, colName := range optionalColumns {
		colName = strings.TrimSpace(strings.ToLower(colName))
		col, found := getColumnByName(colName)
		if !found {
			unknownColumns = append(unknownColumns, colName)
			continue
		}

		activeColumns = append(activeColumns, col)
	}

	if len(unknownColumns) > 0 {
		return activeColumns, fmt.Errorf("Unknown columns: %v", strings.Join(unknownColumns, ", "))
	}

	return activeColumns, nil
}

func filterDefaultColumns() []Column {
	var defaultColumns []Column
	for _, col := range allColumns {
		if col.IsDefault {
			defaultColumns = append(defaultColumns, col)
		}
	}

	return defaultColumns
}

func getColumnByName(name string) (Column, bool) {
	for _, col := range allColumns {
		if col.Name == name {
			return col, true
		}
	}

	return Column{}, false
}
