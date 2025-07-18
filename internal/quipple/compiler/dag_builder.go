package compiler

import (
	"fmt"
	"qp/internal/pipeline/filtering"
	"qp/internal/pkgdata"
	"qp/internal/quipple/query"
)

func RunDAG(expr Expr, pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	root, err := BuildFilterDAG(expr)
	if err != nil {
		return nil, err
	}

	input := make(chan *pkgdata.PkgInfo)
	output := root.Stream(input)

	go func() {
		defer close(input)
		for _, pkg := range pkgs {
			input <- pkg
		}
	}()

	var filtered []*pkgdata.PkgInfo
	for pkg := range output {
		filtered = append(filtered, pkg)
	}

	return filtered, nil
}

func BuildFilterDAG(expr Expr) (FilterNode, error) {
	switch expr := expr.(type) {
	case *AndExpr:
		return buildAndNode(expr)
	case *OrExpr:
		return buildOrNode(expr)
	case *NotExpr:
		return buildNotNode(expr)
	case *QueryExpr:
		return buildQueryNode(expr)
	default:
		return nil, fmt.Errorf("unsupported expression type")
	}
}

func buildAndNode(expr *AndExpr) (FilterNode, error) {
	left, err := BuildFilterDAG(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := BuildFilterDAG(expr.Right)
	if err != nil {
		return nil, err
	}

	return &AndNode{Left: left, Right: right}, nil
}

func buildOrNode(expr *OrExpr) (FilterNode, error) {
	left, err := BuildFilterDAG(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := BuildFilterDAG(expr.Right)
	if err != nil {
		return nil, err
	}

	return &OrNode{left, right}, nil
}

func buildNotNode(expr *NotExpr) (FilterNode, error) {
	inner, err := BuildFilterDAG(expr.Inner)
	if err != nil {
		return nil, err
	}

	return &NotNode{inner}, nil
}

func buildQueryNode(expr *QueryExpr) (FilterNode, error) {
	conditions, err := filtering.QueriesToConditions([]query.FieldQuery{expr.Query})
	if err != nil {
		return nil, err
	}

	if len(conditions) != 1 {
		return nil, fmt.Errorf("expected 1 condition, got %d", len(conditions))
	}

	return &QueryNode{Filter: conditions[0].Filter}, nil
}
