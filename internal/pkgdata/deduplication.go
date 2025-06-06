package pkgdata

import (
	"qp/internal/consts"
	"strings"
)

var originPriorityOrder = []string{
	consts.OriginPacman,
	consts.OriginDeb,
	consts.OriginRpm,
	consts.OriginOpkg,
	consts.OriginBrew,
	consts.OriginPipx,
	consts.OriginNpm,
}

const (
	node         = "node-"
	nodejs       = "nodejs-"
	python       = "python-"
	python3      = "python3-"
	pythonNoDash = "python"
)

func Deduplicate(pkgs []*PkgInfo) []*PkgInfo {
	nameGroups := make(map[string][]*PkgInfo)

	for _, pkg := range pkgs {
		normalizedName := normalizeName(pkg.Name, pkg.Origin)
		nameGroups[normalizedName] = append(nameGroups[normalizedName], pkg)
	}

	var result []*PkgInfo

	for _, group := range nameGroups {
		if len(group) == 1 {
			result = append(result, group[0])
			continue
		}

		bestPkg := bestPriority(group)
		result = append(result, bestPkg)
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

func bestPriority(pkgs []*PkgInfo) *PkgInfo {
	if len(pkgs) == 1 {
		return pkgs[0]
	}

	priorityMap := make(map[string]int)
	for i, origin := range originPriorityOrder {
		priorityMap[origin] = i
	}

	bestPkg := pkgs[0]
	bestPriority := priorityMap[bestPkg.Origin]

	for _, pkg := range pkgs[1:] {
		priority := priorityMap[pkg.Origin]
		if priority < bestPriority {
			bestPkg = pkg
			bestPriority = priority
		}
	}

	return bestPkg
}
