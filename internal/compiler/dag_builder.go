package compiler

import (
	"fmt"
	"qp/internal/ast"
	"qp/internal/pipeline/filtering"
	"qp/internal/pkgdata"
	"qp/internal/query"
)

func RunDAG(expr ast.Expr, pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
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

func BuildFilterDAG(expr ast.Expr) (FilterNode, error) {
	switch expr := expr.(type) {
	case *ast.AndExpr:
		return buildAndNode(expr)
	case *ast.OrExpr:
		return buildOrNode(expr)
	case *ast.NotExpr:
		return buildNotNode(expr)
	case *ast.QueryExpr:
		return buildQueryNode(expr)
	default:
		return nil, fmt.Errorf("unsupported expression type")
	}
}

func buildAndNode(expr *ast.AndExpr) (FilterNode, error) {
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

func buildOrNode(expr *ast.OrExpr) (FilterNode, error) {
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

func buildNotNode(expr *ast.NotExpr) (FilterNode, error) {
	inner, err := BuildFilterDAG(expr.Inner)
	if err != nil {
		return nil, err
	}

	return &NotNode{inner}, nil
}

func buildQueryNode(expr *ast.QueryExpr) (FilterNode, error) {
	conditions, err := filtering.QueriesToConditions([]query.FieldQuery{expr.Query})
	if err != nil {
		return nil, err
	}

	if len(conditions) != 1 {
		return nil, fmt.Errorf("expected 1 condition, got %d", len(conditions))
	}

	return &QueryNode{Filter: conditions[0].Filter}, nil
}
