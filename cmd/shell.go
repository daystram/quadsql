package cmd

import (
	"fmt"
	"os"

	"github.com/daystram/quadsql/handlers"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func interactive(cmd *cobra.Command, args []string) {
	var err error
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
		err = handlers.HandleQuery(result)
	}

	switch err {
	case nil, promptui.ErrInterrupt:
		os.Exit(0)
	default:
		fmt.Printf("an error has occurred. %v\n", err)
		os.Exit(1)
	}
}
