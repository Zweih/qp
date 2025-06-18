package main

import (
	"errors"
	"qp/internal/config"
	"qp/internal/consts"
	out "qp/internal/display"
	"qp/internal/pkgdata"
	"qp/internal/quipple/syntax"
)

func globalPackageSort(
	allPkgs []*pkgdata.PkgInfo,
	cfg *config.Config,
) ([]*pkgdata.PkgInfo, error) {
	comparator, err := pkgdata.GetComparator(cfg.SortOption.Field, cfg.SortOption.Asc)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	if len(allPkgs) >= pkgdata.ConcurrentSortThreshold {
		return pkgdata.SortConcurrently(allPkgs, comparator), nil
	}

	return pkgdata.SortNormally(allPkgs, comparator), nil
}

func trimPackagesLen(
	pkgs []*pkgdata.PkgInfo,
	cfg *config.Config,
) []*pkgdata.PkgInfo {
	if cfg.Limit < 1 || len(pkgs) <= cfg.Limit {
		return pkgs
	}

	switch cfg.LimitMode {
	case syntax.LimitEnd:
		return pkgs[:cfg.Limit]
	case syntax.LimitMid:
		start := (len(pkgs) - cfg.Limit) / 2
		end := start + cfg.Limit
		return pkgs[start:end]
	case syntax.LimitStart:
		fallthrough
	default:
		cutoffIdx := len(pkgs) - cfg.Limit
		return pkgs[cutoffIdx:]
	}
}

func renderOutput(pkgs []*pkgdata.PkgInfo, cfg *config.Config) error {
	switch cfg.OutputFormat {
	case consts.OutputTable:
		out.RenderTable(pkgs, cfg.Fields, cfg.ShowFullTimestamp, cfg.HasNoHeaders)
	case consts.OutputKeyValue:
		out.RenderKeyValue(pkgs, cfg.Fields)
	case consts.OutputJSON:
		out.RenderJSON(pkgs, cfg.Fields)
	default:
		return errors.New("invalid output format")
	}

	return nil
}
