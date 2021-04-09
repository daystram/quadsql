package cmd

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/daystram/quadsql/db"
	"github.com/daystram/quadsql/handlers"
)

var (
	genSeed int64
)

var generateCmd = &cobra.Command{
	Use:   "generate [dimension] [size] [sorted]",
	Short: "Dataset generator",
	Long:  `Generates n-dimensional spatial dataset.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		dimension, _ := strconv.Atoi(args[0])
		size, _ := strconv.Atoi(args[1])
		max := float64(db.MAX_RANGE)
		sorted, _ := strconv.ParseBool(args[2])
		handlers.Generate(DBFile, genSeed, dimension, size, max, sorted)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	f := generateCmd.Flags()
	f.Int64VarP(&genSeed, "seed", "s", 0, "dataset generator seed")
}
