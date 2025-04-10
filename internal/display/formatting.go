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
		if rel.ProviderName != "" {
			virtualFormat = fmt.Sprintf(" (provided by %s)", rel.ProviderName)
		}

		op := relationOpToString(rel.Operator)
		relationOutputs = append(relationOutputs, fmt.Sprintf("%s%s%s%s", rel.Name, op, rel.Version, virtualFormat))
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

func formatDate(pkg *pkgdata.PkgInfo, ctx tableContext) string {
	timestamp := time.Unix(pkg.Timestamp, 0)
	return timestamp.Format(ctx.DateFormat)
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
