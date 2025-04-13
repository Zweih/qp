package display

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"sort"
	"strings"
	"time"
)

func formatRelations(relations []pkgdata.Relation) string {
	if len(relations) == 0 {
		return "-"
	}

	return strings.Join(flattenRelations(relations), ", ")
}

func flattenRelations(relations []pkgdata.Relation) []string {
	relationsAtDepth := pkgdata.GetRelationsByDepth(relations, 1)
	relationOutputs := make([]string, 0, len(relationsAtDepth))

	sort.Slice(relationsAtDepth, func(a int, b int) bool {
		return relationsAtDepth[a].Name < relationsAtDepth[b].Name
	})

	for _, rel := range relationsAtDepth {
		var virtualFormat string
		var whyFormat string

		if rel.ProviderName != "" {
			virtualFormat = fmt.Sprintf(" → %s", rel.ProviderName)
		}

		if rel.Why != "" {
			whyFormat = fmt.Sprintf(" (%s)", rel.Why)
		}

		op := relationOpToString(rel.Operator)
		relationOutputs = append(relationOutputs, fmt.Sprintf("%s%s%s%s%s", rel.Name, op, rel.Version, virtualFormat, whyFormat))
	}

	return relationOutputs
}

func relationOpToString(op pkgdata.RelationOp) string {
	switch op {
	case pkgdata.OpEqual:
		return "="
	case pkgdata.OpLess:
		return "<"
	case pkgdata.OpLessEqual:
		return "<="
	case pkgdata.OpGreater:
		return ">"
	case pkgdata.OpGreaterEqual:
		return ">="
	default:
		return ""
	}
}

func pkgTypeToString(pkgType pkgdata.PkgType) string {
	switch pkgType {
	case pkgdata.PkgTypePkg:
		return "pkg"
	case pkgdata.PkgTypeSplit:
		return "split"
	case pkgdata.PkgTypeSrc:
		return "src"
	case pkgdata.PkgTypeDebug:
		return "debug"
	default:
		return ""
	}
}

func formatDate(timestamp int64, ctx tableContext) string {
	unixTimestamp := time.Unix(timestamp, 0)
	return unixTimestamp.Format(ctx.DateFormat)
}

func formatSize(size int64) string {
	switch {
	case size >= consts.GB:
		return fmt.Sprintf("%.2f GB", float64(size)/(consts.GB))
	case size >= consts.MB:
		return fmt.Sprintf("%.2f MB", float64(size)/(consts.MB))
	case size >= consts.KB:
		return fmt.Sprintf("%.2f KB", float64(size)/(consts.KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
