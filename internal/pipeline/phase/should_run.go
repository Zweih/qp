package phase

import (
	"qp/internal/config"
	"qp/internal/consts"
	"qp/internal/pipeline/meta"
	"slices"
)

func ShouldAlwaysRun(_ config.Config, _ *meta.PipelineContext) bool {
	return true
}

func ShouldRunFetch(_ config.Config, ctx *meta.PipelineContext) bool {
	return !ctx.UsedCache
}

func ShouldRunReverseDeps(cfg config.Config, ctx *meta.PipelineContext) bool {
	if !ctx.UsedCache {
		return false
	}

	_, hasFilter := cfg.FilterQueries[consts.FieldRequiredBy]
	hasField := slices.Contains(cfg.Fields, consts.FieldRequiredBy)

	return hasField || hasFilter
}

func ShouldRunSaveCache(_ config.Config, ctx *meta.PipelineContext) bool {
	return !ctx.UsedCache
}

func ShouldRunFiltering(cfg config.Config, _ *meta.PipelineContext) bool {
	return len(cfg.FilterQueries) > 0
}
