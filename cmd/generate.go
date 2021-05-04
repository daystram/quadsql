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
	Use:   "generate [dimension] [size] [distribution] [sorted]",
	Short: "Dataset generator",
	Long:  `Generates n-dimensional spatial dataset.`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		dimension, _ := strconv.Atoi(args[0])
		size, _ := strconv.Atoi(args[1])
		max := float64(db.MAX_RANGE)
		distribution := args[2]
		sorted, _ := strconv.ParseBool(args[3])
		handlers.Generate(DBFile, genSeed, dimension, size, max, distribution, sorted)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	f := generateCmd.Flags()
	f.Int64VarP(&genSeed, "seed", "s", 0, "dataset generator seed")
}
