package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/daystram/quadsql/db"
	"github.com/daystram/quadsql/handlers"
)

func shell(source string, buildIndex, isPointQuad bool) (err error) {
	var database db.DB
	fmt.Printf("Source DB: %s\n", source)
	if database, err = db.OpenDB(source); err != nil {
		return
	}
	defer database.Close()
	h := handlers.InitHandlers(&database, &handlers.QueryConfig{
		UseIndex: buildIndex,
	})

	if buildIndex {
		if err = h.BuildIndex(isPointQuad); err != nil {
			return
		}
	}

	fmt.Println("Use /help to show help page")
	fmt.Println("Use /exit or Ctrl+C to exit")

	var result string
	prompt := promptui.Prompt{
		Label: ">",
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		},
	}

	for err == nil {
		if result, err = prompt.Run(); err != nil {
			break
		}
		err = h.HandleCommand(result)
	}

	return
}
