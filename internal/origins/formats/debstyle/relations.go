package debstyle

import (
	"qp/internal/pkgdata"
	"strings"
)

func ParseRelations(relStr string) []pkgdata.Relation {
	var rels []pkgdata.Relation

	for entry := range strings.SplitSeq(relStr, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		for alt := range strings.SplitSeq(entry, "|") {
			alt = strings.TrimSpace(alt)

			rel := pkgdata.Relation{
				Depth: 1,
			}

			if i := strings.Index(alt, " ("); i != -1 {
				rel.Name = strings.TrimSpace(alt[:i])
				versionPart := strings.TrimSpace(alt[i+1:])
				versionPart = strings.Trim(versionPart, "() ")

				tokens := strings.SplitN(versionPart, " ", 2)
				if len(tokens) == 2 {
					rel.Operator = pkgdata.StringToOperator(tokens[0])
					rel.Version = tokens[1]
				}
			} else {
				rel.Name = alt
			}

			rels = append(rels, rel)
		}
	}

	return rels
}
