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
	Count             int
	AllPackages       bool
	ShowHelp          bool
	ShowVersion       bool
	OutputJson        bool
	HasNoHeaders      bool
	ShowFullTimestamp bool
	DisableProgress   bool
	NoCache           bool
	RegenCache        bool
	SortOption        syntax.SortOption
	Fields            []consts.FieldType
	FieldQueries      []syntax.FieldQuery
}

type ConfigProvider interface {
	GetConfig() (*Config, error)
}

type CliConfigProvider struct{}
