package rpm

import (
	"qp/internal/pkgdata"
)

type ExprNode interface {
	Evaluate(installedPkgs map[string]*pkgdata.PkgInfo) []pkgdata.Relation
}

type PkgExpr struct {
	Name     string
	Version  string
	Operator pkgdata.RelationOp
}

type AndExpr struct {
	Left, Right ExprNode
}

type OrExpr struct {
	Left, Right ExprNode
}

type CondExpr struct {
	MainExpr ExprNode
	Cond     ExprNode
	IsUnless bool
}

func (pExpr *PkgExpr) Evaluate(_ map[string]*pkgdata.PkgInfo) []pkgdata.Relation {
	return []pkgdata.Relation{{
		Name:      pExpr.Name,
		Version:   pExpr.Version,
		Operator:  pExpr.Operator,
		Depth:     1,
		IsComplex: false,
	}}
}

func (a *AndExpr) Evaluate(installedMap map[string]*pkgdata.PkgInfo) []pkgdata.Relation {
	var result []pkgdata.Relation
	result = append(result, a.Left.Evaluate(installedMap)...)
	result = append(result, a.Right.Evaluate(installedMap)...)
	return result
}

func (o *OrExpr) Evaluate(installedMap map[string]*pkgdata.PkgInfo) []pkgdata.Relation {
	var result []pkgdata.Relation
	result = append(result, o.Left.Evaluate(installedMap)...)
	result = append(result, o.Right.Evaluate(installedMap)...)
	return result
}

func (cond *CondExpr) Evaluate(installedMap map[string]*pkgdata.PkgInfo) []pkgdata.Relation {
	conditionMet := cond.evaluateCond(installedMap)

	if cond.IsUnless {
		if !conditionMet {
			return cond.MainExpr.Evaluate(installedMap)
		}
	} else {
		if conditionMet {
			return cond.MainExpr.Evaluate(installedMap)
		}
	}

	return nil
}

func (cond *CondExpr) evaluateCond(installedMap map[string]*pkgdata.PkgInfo) bool {
	conditionRels := cond.Cond.Evaluate(installedMap)
	for _, rel := range conditionRels {
		installedPkg, exists := installedMap[rel.Name]
		if !exists {
			continue
		}

		if rel.Version != "" {
			if !versionSatisfies(installedPkg.Version, rel.Operator, rel.Version) {
				continue
			}
		}

		return true
	}

	return false
}

func versionSatisfies(
	installedVersion string,
	operator pkgdata.RelationOp,
	requiredVersion string,
) bool {
	comparison := compareVersions(installedVersion, requiredVersion)

	switch operator {
	case pkgdata.OpEqual:
		return comparison == 0
	case pkgdata.OpGreater:
		return comparison > 0
	case pkgdata.OpGreaterEqual:
		return comparison >= 0
	case pkgdata.OpLess:
		return comparison < 0
	case pkgdata.OpLessEqual:
		return comparison <= 0
	case pkgdata.OpNone:
		return true
	default:
		return false
	}
}
