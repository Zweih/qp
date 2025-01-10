package pkgdata

import "sort"

func SortPackages(pkgs []PackageInfo, sortBy string) {
	switch sortBy {
	case "alpabetical":
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Name < pkgs[j].Name
		})

	default: // date is the default sort
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Timestamp.After(pkgs[j].Timestamp)
		})
	}
}
