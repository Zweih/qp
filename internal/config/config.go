package config

import (
	"qp/internal/consts"
	"qp/internal/quipple/compiler"
	"qp/internal/quipple/query"
	"qp/internal/quipple/syntax"
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
	ShowCompletion    string
	CacheOnly         string
	CacheWorker       string
	OutputFormat      string
	LimitMode         syntax.LimitMode
	SortOption        syntax.SortOption
	Fields            []consts.FieldType
	FieldQueries      []query.FieldQuery
	QueryExpr         compiler.Expr
}

type ConfigProvider interface {
	GetConfig() (*Config, error)
}

type CliConfigProvider struct{}
