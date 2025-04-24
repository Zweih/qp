package ast

import (
	"qp/internal/pkgdata"
	"qp/internal/query"
)

type ExprType int32

const (
	AndExprType ExprType = iota
	OrExprType
	NotExprType
	QueryExprType
)

type Expr interface {
	Eval(pkg *pkgdata.PkgInfo) bool
}

type AndExpr struct {
	Left, Right Expr
}

type OrExpr struct {
	Left, Right Expr
}

type NotExpr struct {
	Inner Expr
}

type QueryExpr struct {
	Query query.FieldQuery
}

func (a *AndExpr) Eval(pkg *pkgdata.PkgInfo) bool {
	return a.Left.Eval(pkg) && a.Right.Eval(pkg)
}

func (o *OrExpr) Eval(pkg *pkgdata.PkgInfo) bool {
	return o.Left.Eval(pkg) || o.Right.Eval(pkg)
}

func (n *NotExpr) Eval(pkg *pkgdata.PkgInfo) bool {
	return !n.Inner.Eval(pkg)
}

func (q *QueryExpr) Eval(pkg *pkgdata.PkgInfo) bool {
	panic("QueryExpr.Eval should not be used directly before compiling to FilterConditions")
}
