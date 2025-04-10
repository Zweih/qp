package phase

import (
	"fmt"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/pipeline/filtering"
	"qp/internal/pipeline/meta"
	"qp/internal/pkgdata"
)

func LoadCacheStep(
	cfg *config.Config,
	_ []*PkgInfo,
	_ ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if cfg.NoCache || cfg.RegenCache {
		return nil, nil
	}

	pkgPtrs, err := pkgdata.LoadProtoCache(pipelineCtx.CachePath)
	if err == nil {
		pipelineCtx.UsedCache = true
	}

	// TODO: use ProgressReporter to report cache status
	return pkgPtrs, nil
}

// TODO: add progress reporting
func FetchStep(
	_ *config.Config,
	pkgPtrs []*PkgInfo,
	_ ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if !pipelineCtx.UsedCache {
		var err error
		pkgPtrs, err = pkgdata.FetchPackages()
		if err != nil {
			out.WriteLine(fmt.Sprintf(
				"Warning: Some packages may be missing due to corrupted package database: %v",
				err,
			))
		}
	}

	return pkgPtrs, nil
}

func ResolveDepTreeStep(
	_ *config.Config,
	pkgPtrs []*pkgdata.PkgInfo,
	reportProgress ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if pipelineCtx.UsedCache {
		return pkgPtrs, nil
	}

	return pkgdata.ResolveDependencyTree(pkgPtrs, reportProgress)
}

// TODO: add progress reporting
func SaveCacheStep(
	cfg *config.Config,
	pkgPtrs []*PkgInfo,
	_ ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if !cfg.NoCache && !pipelineCtx.UsedCache {
		// TODO: we can probably save the file concurrently
		err := pkgdata.SaveProtoCache(pkgPtrs, pipelineCtx.CachePath)
		if err != nil {
			out.WriteLine(fmt.Sprintf("Warning: Error saving cache: %v", err))
		}
	}

	return pkgPtrs, nil
}

func FilterStep(
	cfg *config.Config,
	pkgPtrs []*PkgInfo,
	reportProgress ProgressReporter,
	_ *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if len(cfg.FieldQueries) == 0 {
		return pkgPtrs, nil
	}

	filterConditions, err := filtering.QueriesToConditions(cfg.FieldQueries)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	return pkgdata.FilterPackages(pkgPtrs, filterConditions, reportProgress), nil
}

func SortStep(
	cfg *config.Config,
	pkgPtrs []*PkgInfo,
	reportProgress ProgressReporter,
	_ *meta.PipelineContext,
) ([]*PkgInfo, error) {
	comparator, err := pkgdata.GetComparator(cfg.SortOption.Field, cfg.SortOption.Asc)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	phase := "Sorting packages"

	// threshold is 500 as that is where merge sorting chunk performance overtakes timsort
	if len(pkgPtrs) < pkgdata.ConcurrentSortThreshold {
		return pkgdata.SortNormally(pkgPtrs, comparator, phase, reportProgress), nil
	}

	return pkgdata.SortConcurrently(pkgPtrs, comparator, phase, reportProgress), nil
}
