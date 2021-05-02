package handlers

import (
	"fmt"

	"github.com/daystram/quadsql/data"
	"github.com/daystram/quadsql/db"
	"github.com/manifoldco/promptui"
)

type QueryConfig struct {
	ShowTime    bool
	UseIndex    bool
	IndexReady  bool
	IsPointQuad bool // true: Point; false: Region
}

type Handler struct {
	database     *db.DB
	config       *QueryConfig
	index        *data.QuadNode
	lastExecTime float64
}

func InitHandlers(database *db.DB, config *QueryConfig) Handler {
	return Handler{
		database: database,
		config:   config,
	}
}

func (h *Handler) HandleCommand(command string) (err error) {
	switch command {
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
		fmt.Printf("DB Row Count    : %d\n", len(h.database.Table))
		fmt.Printf("Index Enabled   : %s\n", map[bool]string{true: "Yes", false: "No"}[h.config.UseIndex])
		fmt.Printf("Index Ready     : %s\n", map[bool]string{true: "Yes", false: "No"}[h.config.IndexReady])
		fmt.Printf("Index Nodes     : %d nodes\n", count)
		fmt.Printf("Index Max Depth : %d\n", depth)
		fmt.Printf("Last Exec Time  : %.3f µs (%.3f ms)\n", h.lastExecTime/1e3, h.lastExecTime/1e6)
	case "":
		break
	default:
		err = h.performQuery(command)
		if h.config.ShowTime {
			fmt.Printf("Exec time: %.3f µs (%.3f ms)\n", h.lastExecTime/1e3, h.lastExecTime/1e6)
		}
		if err == ErrInvalidQuery {
			err = nil
			fmt.Println("Invalid query statement!")
		}
	}
	return
}
