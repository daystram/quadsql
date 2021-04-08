package handlers

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"

	"github.com/daystram/quadsql/data"
	"github.com/daystram/quadsql/db"
)

type QueryConfig struct {
	ShowTime   bool
	UseIndex   bool
	IndexReady bool
}

type Handler struct {
	database     *db.DB
	config       *QueryConfig
	index        *data.QuadNode
	lastExecTime int64
}

func InitHandlers(database *db.DB, config *QueryConfig) Handler {
	return Handler{
		database: database,
		config:   config,
	}
}

func (h *Handler) BuildIndex() (err error) {
	fmt.Printf("Building index... ")
	// TODO: build index
	h.index = &data.QuadNode{
		Point: data.Point{
			Coordinate: []float32{1, 2},
		},
		Children: []*data.QuadNode{{
			Point: data.Point{
				Coordinate: []float32{1, 2},
			},
		}, {
			Point: data.Point{
				Coordinate: []float32{1, 2},
			},
		}, {
			Point: data.Point{
				Coordinate: []float32{1, 2},
			},
		}, {
			Point: data.Point{
				Coordinate: []float32{1, 2},
			},
		}},
	}
	h.config.IndexReady = true
	fmt.Println("Done")
	return
}

func (h *Handler) HandleQuery(query string) (err error) {
	switch query {
	case "/exit":
		err = promptui.ErrInterrupt
	case "/index on":
		if h.config.IndexReady {
			h.config.UseIndex = true
			fmt.Println("Index enabled")
		} else {
			fmt.Println("Index has not been initialized, use '/index rebuild'")
		}
	case "/index off":
		h.config.UseIndex = false
		fmt.Println("Index disabled")
	case "/index rebuild":
		err = h.BuildIndex()
	case "/time":
		if h.config.ShowTime {
			h.config.ShowTime = false
			fmt.Println("Time report disabled")
		} else {
			h.config.ShowTime = true
			fmt.Println("Time report enabled")
		}
	case "/info":
		count, depth := countNodes(h.index)
		fmt.Printf("Dimension       : %d\n", h.database.Dimension)
		fmt.Printf("DB Row Count    : %d\n", h.database.RowCount)
		fmt.Printf("Index Enabled   : %s\n", map[bool]string{true: "Yes", false: "No"}[h.config.UseIndex])
		fmt.Printf("Index Ready     : %s\n", map[bool]string{true: "Yes", false: "No"}[h.config.IndexReady])
		fmt.Printf("Index Nodes     : %d nodes\n", count)
		fmt.Printf("Index Max Depth : %d\n", depth)
		fmt.Printf("Last Exec Time  : %d ms\n", h.lastExecTime/1000)
	default:
		err = h.performQuery(query)
		if h.config.ShowTime {
			fmt.Printf("Exec time: %d ms\n", h.lastExecTime/1000)
		}

	}
	return
}

func (h *Handler) performQuery(query string) (err error) {
	start := time.Now()
	fmt.Printf("QUERY: %s \n", query)
	// TODO
	fmt.Printf("%+v\n", h.index)
	h.lastExecTime = time.Since(start).Nanoseconds()
	return
}

func countNodes(node *data.QuadNode) (count, depth int) {
	if node != nil {
		for _, child := range node.Children {
			c, d := countNodes(child)
			count += c
			if d > depth {
				depth = d
			}
		}
		count++
		depth++
	}
	return
}
