package handlers

import (
	"fmt"
	"strconv"
	"strings"

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
	args := append(strings.Split(command, " "), "$")
	switch args[0] {
	case "/exit":
		err = promptui.ErrInterrupt
	case "/svg":
		if h.database.Dimension == 2 {
			var scale float64
			if scale, err = strconv.ParseFloat(args[1], 64); err != nil || scale <= 0 {
				fmt.Println("E: invalid scale, see /help")
				return nil
			}
			if len(args)-1 != 3 {
				fmt.Println("E: invalid output filename, see /help")
				return
			}
			h.DrawSVG(scale, args[2])
		} else {
			fmt.Println("E: only supported for 2D spatial dimensions")
		}
	case "/index":
		switch args[1] {
		case "on":
			if h.config.IndexReady {
				h.config.UseIndex = true
				fmt.Println("Index enabled")
			} else {
				fmt.Println("Index has not been initialized, use '/index rebuild [point|region]'")
			}
		case "off":
			h.config.UseIndex = false
			fmt.Println("Index disabled")
		case "rebuild":
			switch args[2] {
			case "point":
				err = h.BuildIndex(true)
			case "region":
				err = h.BuildIndex(false)
			default:
				fmt.Println("E: invalid index type, see /help")
			}
		default:
			fmt.Println("E: invalid action, see /help")
		}
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
		fmt.Printf("Dimension       : %dD\n", h.database.Dimension)
		fmt.Printf("DB Row Count    : %d\n", len(h.database.Table))
		if h.config.IndexReady {
			fmt.Printf("Index Status    : %s\n", map[bool]string{true: "Enabled", false: "Disabled"}[h.config.UseIndex])
			fmt.Printf("Index Type      : %s Quad-tree\n", map[bool]string{true: "Point", false: "Region"}[h.config.IsPointQuad])
		} else {
			fmt.Printf("Index Status    : Uninitialized\n")
			fmt.Printf("Index Type      : None\n")
		}
		fmt.Printf("Index Nodes     : %d nodes\n", count)
		fmt.Printf("Index Max Depth : %d\n", depth)
		fmt.Printf("Last Exec Time  : %.3f µs (%.3f ms)\n", h.lastExecTime/1e3, h.lastExecTime/1e6)
	case "/help", "/?":
		fmt.Println("/info        : display DB and index statistics")
		fmt.Println("/index [cmd] : switch index [on], [off], or [rebuild [point|region]]")
		fmt.Println("/time        : toggle execution time report")
		fmt.Println("/svg [s] [f] : draw SVG of the Quad-tree with scale [s] into file [f].svg")
		fmt.Println("/help        : show this help page")
		fmt.Println("/exit        : exit quadsql")
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
