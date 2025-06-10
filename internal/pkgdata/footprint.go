package pkgdata

func calculateFreeable(pkg *PkgInfo, installedMap map[string]*PkgInfo) int64 {
	freeable := pkg.Size

	seenDeps := make(map[string]struct{})
	for _, depRel := range pkg.Depends {
		reverseKey := depRel.ProviderKey()
		if reverseKey == "" {
			reverseKey = depRel.Key()
		}

		depPkg, exists := installedMap[reverseKey]
		if !exists {
			continue
		}

		if _, seen := seenDeps[reverseKey]; seen {
			continue
		}

		seenDeps[reverseKey] = struct{}{}

		directReqCount := 0
		for _, reqRel := range depPkg.RequiredBy {
			if directReqCount > 1 {
				break
			}

			if reqRel.Depth == 1 {
				directReqCount++
			}
		}

		if directReqCount == 1 {
			freeable += depPkg.Size
		}
	}

	return freeable
}
