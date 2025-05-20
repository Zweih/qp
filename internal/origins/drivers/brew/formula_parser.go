package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strings"

	json "github.com/goccy/go-json"
)

type InstallReceipt struct {
	Time               int64  `json:"time"`
	InstalledOnRequest bool   `json:"installed_on_request"`
	BuiltAsBottle      bool   `json:"built_as_bottle"`
	PouredFromBottle   bool   `json:"poured_from_bottle"`
	Arch               string `json:"arch"`

	RuntimeDependencies []struct {
		FullName         string `json:"full_name"`
		PkgVersion       string `json:"pkg_version"`
		DeclaredDirectly bool   `json:"declared_directly"`
	} `json:"runtime_dependencies"`
}

type FormulaMetadata struct {
	Name                    string   `json:"name"`
	Desc                    string   `json:"desc"`
	License                 string   `json:"license"`
	Homepage                string   `json:"homepage"`
	ConflictsWith           []string `json:"conflicts_with"`
	OptionalDependencies    []string `json:"optional_dependencies"`
	RecommendedDependencies []string `json:"recommended_dependencies"`
}

func mergeFormulaMetadata(pkg *pkgdata.PkgInfo, formula *FormulaMetadata) {
	if formula == nil {
		return
	}

	pkg.Description = formula.Desc
	pkg.License = formula.License
	pkg.Url = formula.Homepage
	pkg.Conflicts = parseConflicts(formula.ConflictsWith)
	pkg.OptDepends = parseOptDepends(formula.OptionalDependencies, formula.RecommendedDependencies)
}

func parseInstallReceipt(path string, version string) (*pkgdata.PkgInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read receipt JSON: %v", err)
	}

	var receipt InstallReceipt
	if err := json.Unmarshal(data, &receipt); err != nil {
		return nil, fmt.Errorf("failed to parse receipt JSON: %v", err)
	}

	pkgName, err := getPkgNameFromPath(path)
	if err != nil {
		return nil, err
	}

	pkg := &pkgdata.PkgInfo{
		InstallTimestamp: receipt.Time,
		Name:             pkgName,
		Reason:           inferInstallReason(receipt),
		Version:          version,
		Arch:             receipt.Arch,
		PkgType:          typeFormula,
		Depends:          parseDepends(receipt),
	}

	inferBuildDate(pkg, receipt)

	return pkg, nil
}

func inferBuildDate(pkg *pkgdata.PkgInfo, receipt InstallReceipt) {
	if !receipt.PouredFromBottle {
		pkg.BuildTimestamp = receipt.Time
	}
}

func inferInstallReason(receipt InstallReceipt) string {
	switch {
	case receipt.InstalledOnRequest:
		return consts.ReasonExplicit
	default:
		return consts.ReasonDependency // TODO: perhaps this should be blank
	}
}

func getPkgNameFromPath(path string) (string, error) {
	parts := strings.Split(filepath.Clean(path), string(os.PathSeparator))

	if len(parts) >= 3 {
		return parts[len(parts)-3], nil
	}

	return "", fmt.Errorf("unexpected receipt path format: %s", path)
}

func parseDepends(receipt InstallReceipt) []pkgdata.Relation {
	rels := make([]pkgdata.Relation, 0, len(receipt.RuntimeDependencies))

	for _, dep := range receipt.RuntimeDependencies {
		if !dep.DeclaredDirectly {
			continue
		}

		rels = append(rels, pkgdata.Relation{
			Name:     dep.FullName,
			Version:  dep.PkgVersion,
			Operator: pkgdata.OpEqual,
			Depth:    1,
		})
	}

	return rels
}

func parseConflicts(conflicts []string) []pkgdata.Relation {
	rels := make([]pkgdata.Relation, 0, len(conflicts))

	for _, conflict := range conflicts {
		rels = append(rels, pkgdata.Relation{
			Name:  conflict,
			Depth: 1,
		})
	}

	return rels
}

func parseOptDepends(optDeps []string, recDeps []string) []pkgdata.Relation {
	inputDeps := append(optDeps, recDeps...)
	rels := make([]pkgdata.Relation, 0, len(inputDeps))

	for _, dep := range inputDeps {
		rels = append(rels, pkgdata.Relation{
			Name:  dep,
			Depth: 1,
		})
	}

	return rels
}
