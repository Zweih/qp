package brew

import "qp/internal/pkgdata"

type CaskRelation struct {
	Cask    []string `json:"cask"`
	Formula []string `json:"formula"`
}

type CaskMetadata struct {
	Token         string       `json:"token"`
	Caveats       string       `json:"caveats"`
	Desc          string       `json:"desc"`
	Homepage      string       `json:"homepage"`
	ConflictsWith CaskRelation `json:"conflicts_with"`
	DependsOn     CaskRelation `json:"depends_on"`
}

func mergeCaskMetadata(pkg *pkgdata.PkgInfo, cask *CaskMetadata) {
	if cask == nil {
		return
	}

	pkg.Name = cask.Token
	pkg.Description = cask.Desc
	pkg.Url = cask.Homepage
	pkg.Depends = parseCaskRelations(cask.DependsOn)
}

func parseCaskRelations(caskRelation CaskRelation) []pkgdata.Relation {
	formulaRels := parseRawCaskRels(typeFormula, caskRelation.Formula)
	caskRels := parseRawCaskRels(typeCask, caskRelation.Formula)

	return append(formulaRels, caskRels...)
}

func parseRawCaskRels(pkgType string, rawRels []string) []pkgdata.Relation {
	if len(rawRels) < 1 {
		return []pkgdata.Relation{}
	}

	rels := make([]pkgdata.Relation, 0, len(rawRels))

	for _, rawRel := range rawRels {
		rels = append(rels, pkgdata.Relation{
			Name:    rawRel,
			Depth:   1,
			PkgType: pkgType,
		})
	}

	return rels
}
