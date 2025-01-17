package pkgdata

import "sort"

func SortPackages(pkgs []PackageInfo, sortBy string) {
	switch sortBy {
	case "alphabetical":
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Name < pkgs[j].Name
		})

	case "size:desc":
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Size > pkgs[j].Size
		})

	case "size:asc":
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Size < pkgs[j].Size
		})

	default: // date is the default sort
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Timestamp.After(pkgs[j].Timestamp)
		})
	}
}
