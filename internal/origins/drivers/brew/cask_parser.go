package brew

import (
	"fmt"
	"os"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"

	json "github.com/goccy/go-json"
)

type RuntimeDependency struct {
	FullName         string `json:"full_name"`
	DeclaredDirectly bool   `json:"declared_directly"`
}

type CaskReceipt struct {
	Time               int64  `json:"time"`
	InstalledOnRequest bool   `json:"installed_on_request"`
	Arch               string `json:"arch"`
	Source             struct {
		Version string `json:"version"`
	} `json:"source"`
	RuntimeDependencies map[string][]RuntimeDependency `json:"runtime_dependencies"`
}

type CaskMetadata struct {
	Token         string              `json:"token"`
	Caveats       string              `json:"caveats"`
	Desc          string              `json:"desc"`
	Homepage      string              `json:"homepage"`
	ConflictsWith map[string][]string `json:"conflicts_with"`
}

func parseCaskReceipt(name string, path string) (*pkgdata.PkgInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read receipt JSON: %v", err)
	}

	var receipt CaskReceipt
	if err := json.Unmarshal(data, &receipt); err != nil {
		return nil, fmt.Errorf("failed to parse receipt JSON: %v", err)
	}

	reason := consts.ReasonExplicit
	if !receipt.InstalledOnRequest {
		reason = consts.ReasonDependency
	}

	formulaRels := parseCaskDeps(typeFormula, receipt.RuntimeDependencies[typeFormula])
	caskRels := parseCaskDeps(typeCask, receipt.RuntimeDependencies[typeCask])

	pkg := &pkgdata.PkgInfo{
		Name:            name,
		Version:         receipt.Source.Version,
		UpdateTimestamp: receipt.Time,
		Arch:            receipt.Arch,
		Depends:         append(formulaRels, caskRels...),
		Reason:          reason,
		PkgType:         typeCask,
	}

	creationTime, isReliable, err := shared.GetCreationTime(path)
	if err == nil && isReliable {
		pkg.InstallTimestamp = creationTime
	}

	return pkg, nil
}

func parseCaskDeps(pkgType string, deps []RuntimeDependency) []pkgdata.Relation {
	if len(deps) < 1 {
		return []pkgdata.Relation{}
	}

	rels := make([]pkgdata.Relation, 0, len(deps))

	for _, dep := range deps {
		if dep.DeclaredDirectly {
			rels = append(rels, pkgdata.Relation{
				Name:    dep.FullName,
				Depth:   1,
				PkgType: pkgType,
			})
		}
	}

	return rels
}

func mergeCaskMetadata(pkg *pkgdata.PkgInfo, cask *CaskMetadata) {
	if cask == nil {
		return
	}

	pkg.Name = cask.Token
	pkg.Description = cask.Desc
	pkg.Url = cask.Homepage

	formulaRels := parseRawRels(cask.ConflictsWith[typeFormula])
	caskRels := parseRawRels(cask.ConflictsWith[typeCask])

	pkg.Conflicts = append(formulaRels, caskRels...)
}
