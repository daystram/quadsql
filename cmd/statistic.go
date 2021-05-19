package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daystram/quadsql/db"
	"github.com/daystram/quadsql/handlers"
)

var (
	statRuns int64
	noHead   bool
	mode     string
)

var statisticCmd = &cobra.Command{
	Use:   "statistic",
	Short: "Retrieve all statistics",
	Long:  `Generates all run statistics for each point in dataset. Outputs CSV to STDERR.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var database db.DB
		if database, err = db.OpenDB(DBFile); err != nil {
			return
		}
		defer database.Close()
		h := handlers.InitHandlers(&database, &handlers.QueryConfig{})
		if !noHead {
			fmt.Fprint(os.Stderr, "dim,row,index_type,node_count,max_depth,build_time,avg_exec_time,avg_index_access,avg_table_access,avg_point_comp,runs\n")
		}
		for _, m := range strings.Split(mode, ",") {
			switch m {
			case "none":
				h.HandleCommand("/index off")
				fmt.Println("STAT: sampling index=none")
				fmt.Fprintln(os.Stderr, h.CompileStatistics(sampleStats(database, h), statRuns))
			case "point":
				h.HandleCommand("/index rebuild point")
				h.HandleCommand("/index on")
				fmt.Println("STAT: sampling index=point")
				fmt.Fprintln(os.Stderr, h.CompileStatistics(sampleStats(database, h), statRuns))
			case "region":
				h.HandleCommand("/index rebuild region")
				h.HandleCommand("/index on")
				fmt.Println("STAT: sampling index=region")
				fmt.Fprintln(os.Stderr, h.CompileStatistics(sampleStats(database, h), statRuns))
			}
		}
	},
}

func sampleStats(db db.DB, h handlers.Handler) (stat []handlers.Stat) {
	for i := int64(0); i < statRuns; i++ {
		for _, point := range db.Table {
			res, err := h.Execute(handlers.QueryModel{
				Type: "SELECT",
				Condition: handlers.Condition{
					Field: "position",
					Value: point,
				},
			})
			if len(res) == 0 || res[0].Position.CompareTo(point) != 0 || err != nil {
				panic("E: query error")
			}
			stat = append(stat, h.Statistic)
		}
	}
	return
}

func init() {
	rootCmd.AddCommand(statisticCmd)

	f := statisticCmd.Flags()
	f.Int64Var(&statRuns, "runs", 1, "number of run loops")
	f.BoolVar(&noHead, "no-head", false, "exclude CSV header")

	rf := statisticCmd.Flags()
	rf.StringVar(&mode, "mode", "", "comma seperated values of sampled modes")

	cobra.MarkFlagRequired(rf, "mode")
}
