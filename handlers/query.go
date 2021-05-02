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
				Position: []float64{278.99222549719235, 159.8223863234124},
			},
		},
	}
	fmt.Printf("QUERY: %+v \n", queryModel)

	start := time.Now()
	var result QueryResult
	if result, err = h.execute(queryModel); err != nil {
		return
	}
	h.lastExecTime = float64(time.Since(start).Nanoseconds())

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
			switch query.Condition.Field {
			case "id":
				// linear scan on table
				id := query.Condition.Value.(int)
				if point, ok := h.database.Table[id]; ok {
					result = append(result, Row{id, point})
				}
			case "position":
				position := query.Condition.Value.(data.Point)
				if h.config.UseIndex {
					// TODO: point query using index on table
					break
				} else {
					// linear scan on table
					for id, point := range h.database.Table {
						if point.CompareTo(position) == 0 {
							result = append(result, Row{id, point})
							break
						}
					}
				}
			}
		}
	case "DELETE":
		break
	case "INSERT":
		break
	default:
		err = ErrInvalidQuery
	}
	return
}
