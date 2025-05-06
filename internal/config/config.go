package config

import (
	"qp/internal/ast"
	"qp/internal/consts"
	"qp/internal/query"
	"qp/internal/syntax"
)

type Config struct {
	Limit             int
	ShowHelp          bool
	ShowVersion       bool
	OutputJSON        bool
	HasNoHeaders      bool
	ShowFullTimestamp bool
	DisableProgress   bool
	NoCache           bool
	RegenCache        bool
	LimitMode         syntax.LimitMode
	SortOption        syntax.SortOption
	Fields            []consts.FieldType
	FieldQueries      []query.FieldQuery
	QueryExpr         ast.Expr
}

type ConfigProvider interface {
	GetConfig() (*Config, error)
}

type CliConfigProvider struct{}
