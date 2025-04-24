package config

import (
	"qp/internal/consts"
	"qp/internal/syntax"
)

const (
	ReasonExplicit   = "explicit"
	ReasonDependency = "dependency"
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
	FieldQueries      []syntax.FieldQuery
	QueryExpr         syntax.Expr
}

type ConfigProvider interface {
	GetConfig() (*Config, error)
}

type CliConfigProvider struct{}
