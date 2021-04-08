package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/daystram/quadsql/db"
	"github.com/daystram/quadsql/handlers"
)

func interactive(source string, buildIndex bool) (err error) {
	var database db.DB
	if database, err = db.OpenDB(source); err != nil {
		return
	}
	defer database.Close()
	fmt.Printf("Source DB: %s\n", source)
	h := handlers.InitHandlers(&database, &handlers.QueryConfig{
		UseIndex: buildIndex,
	})

	if buildIndex {
		if err = h.BuildIndex(); err != nil {
			return
		}
	}

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
