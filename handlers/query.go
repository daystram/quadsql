package handlers

import (
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/daystram/quadsql/db"
)

func HandleQuery(db *db.DB, query string) (err error) {
	switch query {
	case "/exit":
		err = promptui.ErrInterrupt
	default:
		err = performQuery(query)
	}
	return
}

func performQuery(query string) (err error) {
	fmt.Printf("QUERY: %s\n", query)
	return
}
