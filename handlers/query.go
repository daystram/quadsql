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
	Condition Condition
}

type Condition struct {
	Field string // id, position
	Value interface{}
}

type QueryResult struct {
	Count  int
	Result []struct {
		ID    int
		Point data.Point
	}
}

func (h *Handler) performQuery(query string) (err error) {
	start := time.Now()
	fmt.Printf("QUERY: %s \n", query)

	queryModel := QueryModel{
		Type: "SELECT",
		Condition: Condition{
			Field: "position",
			Value: data.Point{
				Position: []float32{1, 2},
			},
		},
	}

	var result QueryResult
	if result, err = h.execute(queryModel); err != nil {
		return
	}

	fmt.Println("  id     position")
	fmt.Println("----------------------------")
	for _, row := range result.Result {
		fmt.Printf("  %d     %s\n", row.ID, row.Point)
	}
	fmt.Printf("\n%d rows\n", result.Count)

	h.lastExecTime = time.Since(start).Nanoseconds()
	return
}

func (h *Handler) execute(query QueryModel) (result QueryResult, err error) {
	switch query.Type {
	case "SELECT":
		break
	case "DELETE":
		break
	case "INSERT":
		break
	default:
		err = ErrInvalidQuery
	}
	return
}
