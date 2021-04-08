package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/daystram/quadsql/db"
	"github.com/daystram/quadsql/handlers"
)

var database db.DB

func interactive(source string) (err error) {
	if database, err = db.OpenDB(source); err != nil {
		return
	}
	defer database.Close()
	fmt.Printf("Source DB: %s\n", source)

	fmt.Printf("Building index... ")
	// TODO: build index
	fmt.Println("Done!")

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
		err = handlers.HandleQuery(&database, result)
	}

	return
}
