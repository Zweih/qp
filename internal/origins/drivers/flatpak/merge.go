package flatpak

import (
	"qp/internal/pkgdata"
)

func mergeExtensions(pkgs []*pkgdata.PkgInfo) []*pkgdata.PkgInfo {
	pkgMap := make(map[string]*pkgdata.PkgInfo)
	results := make([]*pkgdata.PkgInfo, 0, len(pkgs))

	for _, pkg := range pkgs {
		uniqueKey := pkg.Name + "|" + pkg.Env
		pkgMap[uniqueKey] = pkg
	}

	for _, pkg := range pkgs {
		if pkg, exists := pkgMap[pkg.Name+"|"+pkg.Env]; exists {
			optDepRels := []pkgdata.Relation{}

			for _, optDepRel := range pkg.OptDepends {
				if optDepRel.PkgType == pkg.Name {
					optKey := optDepRel.Name + "|" + pkg.Env
					if optDepPkg, exists := pkgMap[optKey]; exists {
						pkg.Size += optDepPkg.Size
					}

					delete(pkgMap, optKey)
					continue
				}

				optDepRels = append(optDepRels, optDepRel)
			}

			pkg.OptDepends = optDepRels
		}
	}

	for _, pkg := range pkgMap {
		results = append(results, pkg)
	}

	return results
}
