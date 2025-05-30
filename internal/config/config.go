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
	HasNoHeaders      bool
	ShowFullTimestamp bool
	DisableProgress   bool
	NoCache           bool
	RegenCache        bool
	CacheOnly         string
	CacheWorker       string
	OutputFormat      string
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
