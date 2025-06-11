package pkgdata

import (
	"qp/internal/consts"
	"sort"
	"strings"
)

const (
	node         = "node-"
	nodejs       = "nodejs-"
	python       = "python-"
	python3      = "python3-"
	pythonNoDash = "python"
)

var originExclusion = map[string]struct{}{
	consts.OriginPipx: {},
}

func EnrichAcrossOrigins(pkgs []*PkgInfo) []*PkgInfo {
	nameGroups := make(map[string][]*PkgInfo)

	for _, pkg := range pkgs {
		normalizedName := normalizeName(pkg.Name, pkg.Origin)
		nameGroups[normalizedName] = append(nameGroups[normalizedName], pkg)
	}

	var result []*PkgInfo

	for _, group := range nameGroups {
		for _, pkg := range group {
			if _, exists := originExclusion[pkg.Origin]; exists {
				continue
			}

			pkgCopy := *pkg

			var otherOrigins []string
			var otherEnvs []string

			for _, otherPkg := range group {
				if pkg == otherPkg {
					continue
				}

				if otherPkg.Origin != pkg.Origin {
					otherOrigins = append(otherOrigins, otherPkg.Origin)
					continue
				}

				if pkg.Env != "" {
					otherEnvs = append(otherEnvs, otherPkg.Env)
				}
			}

			sort.Strings(otherOrigins)
			sort.Strings(otherEnvs)
			pkgCopy.AlsoIn = otherOrigins
			pkgCopy.OtherEnvs = otherEnvs

			result = append(result, &pkgCopy)
		}
	}

	return result
}

func normalizeName(name, origin string) string {
	switch origin {
	case consts.OriginPacman:
		return normalizeArch(name)
	case consts.OriginDeb:
		return normalizeDeb(name)
	case consts.OriginRpm:
		return normalizeRpm(name)
	default:
		return strings.ToLower(name)
	}
}

func normalizeArch(name string) string {
	if strings.HasPrefix(name, python) {
		return strings.TrimPrefix(name, python)
	}

	if strings.HasPrefix(name, nodejs) {
		return strings.TrimPrefix(name, nodejs)
	}

	return strings.ToLower(name)
}

func normalizeDeb(name string) string {
	if strings.HasPrefix(name, python3) {
		return strings.TrimPrefix(name, python3)
	}

	if strings.HasPrefix(name, node) {
		return strings.TrimPrefix(name, node)
	}

	return strings.ToLower(name)
}

func normalizeRpm(name string) string {
	if strings.HasPrefix(name, python3) {
		return strings.TrimPrefix(name, python3)
	}

	// versioned python (python39-xyx), we can't guarantee what version fedora will support next
	if strings.HasPrefix(name, pythonNoDash) && strings.Contains(name, "-") {
		parts := strings.SplitN(name, "-", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], pythonNoDash) {
			return parts[1]
		}
	}

	if strings.HasPrefix(name, nodejs) {
		return strings.TrimPrefix(name, nodejs)
	}

	return strings.ToLower(name)
}
