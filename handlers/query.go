package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/daystram/quadsql/data"
)

/*
	Example queries:
		SELECT *
		SELECT * WHERE position=Point(1,2)
*/

type QueryConfig struct {
	ShowStat    bool
	UseIndex    bool
	IndexReady  bool
	IsPointQuad bool // true: Point; false: Region
}

type QueryModel struct {
	Type      string // SELECT, DELETE, INSERT
	Condition Condition
}

type Condition struct {
	Field string // id, position
	Value interface{}
}

type QueryResult []Row

type Row struct {
	ID       int
	Position data.Point
}

func (h *Handler) performQuery(query string) (err error) {

	// query parser
	queryModel := QueryModel{Type: "SELECT"}
	condition := strings.TrimPrefix(strings.TrimPrefix(strings.ToLower(query), "select *"), " where ")
	switch args := strings.Split(condition, "="); args[0] {
	case "id":
		var id int
		if id, err = strconv.Atoi(args[1]); err != nil {
			return ErrInvalidQuery
		}
		queryModel.Condition = Condition{
			Field: "id",
			Value: id,
		}
	case "position":
		var position data.Point
		if position, err = data.ParsePoint(args[1]); err != nil {
			return ErrInvalidQuery
		}
		queryModel.Condition = Condition{
			Field: "position",
			Value: position,
		}
	case "":
	default:
		return ErrInvalidQuery
	}

	var result QueryResult
	if result, err = h.execute(queryModel); err != nil {
		return
	}

	fmt.Println("\n    id\t position")
	fmt.Println("----------------------------")
	for _, row := range result {
		fmt.Printf("  %d\t %s\n", row.ID, row.Position)
	}
	fmt.Printf("\nMatched %d row%s\n", len(result), map[bool]string{true: "", false: "s"}[len(result) == 1])

	return
}

func (h *Handler) execute(query QueryModel) (result QueryResult, err error) {
	start := time.Now()
	accessIndex, accessTable, comparePoint := 0, 0, 0
	switch query.Type {
	case "SELECT":
		if query.Condition.Field == "" {
			for id, point := range h.database.Table {
				accessTable++
				result = append(result, Row{id, point})
			}
		} else {
			switch query.Condition.Field {
			case "id":
				// direct table access
				id := query.Condition.Value.(int)
				if id < len(h.database.Table) {
					accessTable++
					result = append(result, Row{id, h.database.Table[id]})
				}
			case "position":
				position := query.Condition.Value.(data.Point)
				if h.config.UseIndex {
					if h.config.IsPointQuad {
						// point query using Point index on table
						node := h.index
						for node != nil {
							accessIndex++
							comparePoint++
							if node.Centre.CompareTo(position) == 0 {
								accessTable++
								result = append(result, Row{*node.PointID, h.database.Table[*node.PointID]})
								break
							} else {
								node = node.Children[getQuadrant(node.Centre, position)]
							}
						}
					} else {
						// point query using Region index on table
						accessIndex++
						node := h.index
						for node.PointID == nil {
							accessIndex++
							node = node.Children[getQuadrant(node.Centre, position)]
						}
						comparePoint++
						if node.Centre.CompareTo(position) == 0 {
							accessTable++
							result = append(result, Row{*node.PointID, h.database.Table[*node.PointID]})
						}
					}
				} else {
					// linear scan on table
					for id, point := range h.database.Table {
						comparePoint++
						accessTable++
						if point.CompareTo(position) == 0 {
							result = append(result, Row{id, point})
							break
						}
					}
				}
			}
		}
		h.statistic = Stat{
			TimeExec:     float64(time.Since(start).Nanoseconds()),
			AccessIndex:  accessIndex,
			AccessTable:  accessTable,
			ComparePoint: comparePoint,
		}
	default:
		err = ErrInvalidQuery
	}
	return
}
