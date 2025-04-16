package config

import (
	"qp/internal/consts"
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
	SortOption        SortOption
	Fields            []consts.FieldType
	FieldQueries      []FieldQuery
}

type FieldQuery struct {
	IsExistence bool
	Negate      bool
	Field       consts.FieldType
	Match       consts.MatchType
	Depth       int32
	Target      string
}

type SortOption struct {
	Field consts.FieldType
	Asc   bool
}

type ConfigProvider interface {
	GetConfig() (*Config, error)
}

type CliConfigProvider struct{}
