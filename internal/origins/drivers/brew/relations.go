package brew

import "qp/internal/pkgdata"

func parseRawRels(rawRels []string) []pkgdata.Relation {
	if len(rawRels) < 1 {
		return []pkgdata.Relation{}
	}

	rels := make([]pkgdata.Relation, 0, len(rawRels))

	for _, rawRel := range rawRels {
		rels = append(rels, pkgdata.Relation{
			Name:  rawRel,
			Depth: 1,
		})
	}

	return rels
}
