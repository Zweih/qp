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
	OutputJson        bool
	HasNoHeaders      bool
	ShowFullTimestamp bool
	DisableProgress   bool
	NoCache           bool
	SortOption        SortOption
	Fields            []consts.FieldType
	FieldQueries      FieldQueries
}

type (
	FieldQueries    map[consts.FieldType]SubfieldQueries
	SubfieldQueries map[consts.SubfieldType]string
)

type SortOption struct {
	Field consts.FieldType
	Asc   bool
}

type ConfigProvider interface {
	GetConfig() (Config, error)
}

type CliConfigProvider struct{}
