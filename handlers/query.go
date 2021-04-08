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
		DELETE WHERE id = 4
		DELETE WHERE position = Point(1,2)
		INSERT Point(5,6)
*/

type QueryModel struct {
	Type      string // SELECT, DELETE, INSERT
	Condition *Condition
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
	start := time.Now()
	fmt.Printf("QUERY: %s \n", query)

	queryModel := QueryModel{
		Type: "SELECT",
		Condition: &Condition{
			Field: "position",
			Value: data.Point{
				Position: []float64{2, 3},
			},
		},
	}

	var result QueryResult
	if result, err = h.execute(queryModel); err != nil {
		return
	}

	fmt.Println("  id\t position")
	fmt.Println("----------------------------")
	for _, row := range result {
		fmt.Printf("  %d\t %s\n", row.ID, row.Position)
	}
	fmt.Printf("\n%d rows\n", len(result))

	h.lastExecTime = time.Since(start).Nanoseconds()
	return
}

func (h *Handler) execute(query QueryModel) (result QueryResult, err error) {
	switch query.Type {
	case "SELECT":
		if query.Condition == nil {
			for id, point := range h.database.Table {
				result = append(result, Row{id, point})
			}
		} else {
			switch query.Condition.Field {
			case "id":
				id := query.Condition.Value.(int)
				if point, ok := h.database.Table[id]; ok {
					result = append(result, Row{id, point})
				}
			case "position":
				position := query.Condition.Value.(data.Point)
				if h.config.UseIndex {
					break
				} else {
					for id, point := range h.database.Table {
						if point.Equals(position) {
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
