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

func GetColumnJsonValues(pkg pkgdata.PackageInfo, columnNames []string) pkgdata.PackageInfoJson {
	filteredPackage := pkgdata.PackageInfoJson{}

	for _, columnName := range columnNames {
		switch columnName {
		case consts.Date:
			filteredPackage.Timestamp = &pkg.Timestamp
		case consts.Name:
			filteredPackage.Name = pkg.Name
		case consts.Reason:
			filteredPackage.Reason = pkg.Reason
		case consts.Size:
			filteredPackage.Size = pkg.Size // return in bytes for json
		case consts.Version:
			filteredPackage.Version = pkg.Version
		case consts.Depends:
			filteredPackage.Depends = pkg.Depends
		case consts.RequiredBy:
			filteredPackage.RequiredBy = pkg.RequiredBy
		case consts.Provides:
			filteredPackage.Provides = pkg.Provides
		}
	}

	return filteredPackage
}

func GetColumnTableValue(pkg pkgdata.PackageInfo, columnName string, ctx displayContext) string {
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
func formatDate(pkg pkgdata.PackageInfo, ctx displayContext) string {
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
