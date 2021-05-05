package handlers

import (
	"fmt"

	"gonum.org/v1/gonum/stat"

	"github.com/daystram/quadsql/utils"
)

type Stat struct {
	TimeExec       float64
	TimeIndexBuild float64
	AccessIndex    int
	AccessTable    int
	ComparePoint   int
}

func (h *Handler) CompileStatistics(stats []Stat, runs int64) (report string) {
	dim := h.database.Dimension
	row := len(h.database.Table)
	indexType := "none"
	if h.config.UseIndex {
		if h.config.IsPointQuad {
			indexType = "point"
		} else {
			indexType = "region"
		}
	}
	nodeCount, maxDepth := utils.CountNodes(h.index)
	timeIndexBuild := h.Statistic.TimeIndexBuild / 1e6
	avgTimeExec, avgAccessIndex, avgAccessTable, avgComparePoint := aggregateStatistics(stats)
	report += fmt.Sprintf("%d,%d,%s,%d,%d,%f,%f,%f,%f,%f,%d",
		dim, row, indexType, nodeCount, maxDepth, timeIndexBuild, avgTimeExec, avgAccessIndex, avgAccessTable, avgComparePoint, runs)
	return
}

func aggregateStatistics(stats []Stat) (avgTimeExec, avgAccessIndex, avgAccessTable, avgComparePoint float64) {
	var (
		timeExec     []float64
		accessIndex  []float64
		accessTable  []float64
		comparePoint []float64
	)
	for _, stat := range stats {
		timeExec = append(timeExec, stat.TimeExec/1e6)
		accessIndex = append(accessIndex, float64(stat.AccessIndex))
		accessTable = append(accessTable, float64(stat.AccessTable))
		comparePoint = append(comparePoint, float64(stat.ComparePoint))
	}
	return stat.Mean(timeExec, nil), stat.Mean(accessIndex, nil), stat.Mean(accessTable, nil), stat.Mean(comparePoint, nil)
}
