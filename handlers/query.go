package handlers

import (
	"fmt"

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
	database *db.DB
	config   *QueryConfig
	index    data.QuadNode
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
	h.index = data.QuadNode{}
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
	default:
		err = h.performQuery(query)
	}
	return
}

func (h *Handler) performQuery(query string) (err error) {
	fmt.Printf("QUERY: %s \n", query)

	return
}
