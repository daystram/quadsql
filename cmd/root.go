package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	DBFile   string
	NoIndex  bool
	IsRegion bool
)

var rootCmd = &cobra.Command{
	Use:   "quadsql",
	Short: "A Quad-Tree implementation",
	Long:  "A Quad-Tree implementation for high-dimensional spatial databases.",
	Run: func(cmd *cobra.Command, args []string) {
		err := shell(DBFile, !NoIndex, !IsRegion)
		switch err {
		case nil, promptui.ErrInterrupt:
			os.Exit(0)
		default:
			fmt.Printf("an error has occurred. %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&DBFile, "db", "", "DB source file")
	pf.BoolVar(&NoIndex, "no-index", false, "skip index building on start")
	pf.BoolVar(&IsRegion, "region", false, "set mode to Region Quad-tree index, defaults to Point")

	_ = cobra.MarkFlagRequired(pf, "db")
	_ = cobra.MarkFlagFilename(pf, "db")
}
