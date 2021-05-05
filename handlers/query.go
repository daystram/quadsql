package handlers

import (
	"fmt"
	"time"

	"github.com/daystram/quadsql/data"
)

/*
	Example queries:
		SELECT *
		SELECT * WHERE position = Point(1,2)
		// DELETE WHERE id = 4
		// DELETE WHERE position = Point(1,2)
		// INSERT Point(5,6)
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

	// TODO: query parser

	// queryModel := QueryModel{
	// 	Type:      "SELECT",
	// 	Condition: Condition{},
	// }
	queryModel := QueryModel{
		Type: "SELECT",
		Condition: Condition{
			Field: "position",
			Value: data.Point{
				Position: []float64{220.15633674827117, 634.219107312793},
			},
		},
	}
	fmt.Printf("QUERY: %+v \n", queryModel)

	start := time.Now()
	var result QueryResult
	if result, err = h.execute(queryModel); err != nil {
		return
	}
	h.statistic.TimeExec = float64(time.Since(start).Nanoseconds())

	fmt.Println("\n  id\t position")
	fmt.Println("----------------------------")
	for _, row := range result {
		fmt.Printf("  %d\t %s\n", row.ID, row.Position)
	}
	fmt.Printf("\nMatched %d row%s\n", len(result), map[bool]string{true: "", false: "s"}[len(result) == 1])

	return
}

func (h *Handler) execute(query QueryModel) (result QueryResult, err error) {
	switch query.Type {
	case "SELECT":
		if query.Condition.Field == "" {
			for id, point := range h.database.Table {
				result = append(result, Row{id, point})
			}
		} else {
			accessIndex, accessTable, comparePoint := 0, 0, 0
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
			fmt.Printf("Index accesses  : %d\n", accessIndex)
			fmt.Printf("Table accesses  : %d\n", accessTable)
			fmt.Printf("Point comparison: %d\n", comparePoint)
		}
	default:
		err = ErrInvalidQuery
	}
	return
}
