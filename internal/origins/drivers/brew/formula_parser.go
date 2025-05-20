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

type FormulaReceipt struct {
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

func parseFormulaReceipt(path string, version string) (*pkgdata.PkgInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read receipt JSON: %v", err)
	}

	var receipt FormulaReceipt
	if err := json.Unmarshal(data, &receipt); err != nil {
		return nil, fmt.Errorf("failed to parse receipt JSON: %v", err)
	}

	pkgName, err := getPkgNameFromPath(path)
	if err != nil {
		return nil, err
	}

	reason := consts.ReasonExplicit
	if !receipt.InstalledOnRequest {
		reason = consts.ReasonDependency
	}

	pkg := &pkgdata.PkgInfo{
		InstallTimestamp: receipt.Time,
		Name:             pkgName,
		Reason:           reason,
		Version:          version,
		Arch:             receipt.Arch,
		PkgType:          typeFormula,
		Depends:          parseDepends(receipt),
	}

	inferBuildDate(pkg, receipt)

	return pkg, nil
}

func inferBuildDate(pkg *pkgdata.PkgInfo, receipt FormulaReceipt) {
	if !receipt.PouredFromBottle {
		pkg.BuildTimestamp = receipt.Time
	}
}

func getPkgNameFromPath(path string) (string, error) {
	parts := strings.Split(filepath.Clean(path), string(os.PathSeparator))

	if len(parts) >= 3 {
		return parts[len(parts)-3], nil
	}

	return "", fmt.Errorf("unexpected receipt path format: %s", path)
}

func parseDepends(receipt FormulaReceipt) []pkgdata.Relation {
	rels := make([]pkgdata.Relation, 0, len(receipt.RuntimeDependencies))

	for _, dep := range receipt.RuntimeDependencies {
		if !dep.DeclaredDirectly {
			continue
		}

		rels = append(rels, pkgdata.Relation{
			Name:     dep.FullName,
			Version:  dep.PkgVersion,
			Operator: pkgdata.OpEqual,
			PkgType:  typeFormula,
			Depth:    1,
		})
	}

	return rels
}

func mergeFormulaMetadata(pkg *pkgdata.PkgInfo, formula *FormulaMetadata) {
	if formula == nil {
		return
	}

	pkg.Description = formula.Desc
	pkg.License = formula.License
	pkg.Url = formula.Homepage
	pkg.Conflicts = parseRawRels(formula.ConflictsWith)
	optDeps := append(formula.OptionalDependencies, formula.RecommendedDependencies...)
	pkg.OptDepends = parseRawRels(optDeps)
}
