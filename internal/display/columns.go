package display

import (
	"fmt"
	"strings"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

type displayContext struct {
	DateFormat string
}

type PackageInfo = pkgdata.PackageInfo

type Column struct {
	Header string
	Getter func(pkg PackageInfo, ctx displayContext) string
}

var allColumns = map[string]Column{
	consts.DATE: {"DATE", formatDate},
	consts.NAME: {"NAME", func(pkg PackageInfo, _ displayContext) string {
		return pkg.Name
	}},
	consts.VERSION: {"VERSION", func(pkg PackageInfo, _ displayContext) string {
		return pkg.Version
	}},
	consts.REASON: {"REASON", func(pkg PackageInfo, _ displayContext) string {
		return pkg.Reason
	}},
	consts.SIZE: {"SIZE", func(pkg PackageInfo, _ displayContext) string {
		return formatSize(pkg.Size)
	}},
	consts.DEPENDS: {"DEPENDS", func(pkg PackageInfo, _ displayContext) string {
		return formatDependencies(pkg.Depends)
	}},
}

func formatDate(pkg PackageInfo, ctx displayContext) string {
	return pkg.Timestamp.Format(ctx.DateFormat)
}

func formatDependencies(depends []string) string {
	if len(depends) == 0 {
		return "-"
	}
	return strings.Join(depends, ", ")
}

func GetColumnByName(name string) Column {
	col := allColumns[name]
	return col
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
