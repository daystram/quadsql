package handlers

import (
	"fmt"
	"time"

	"github.com/daystram/quadsql/data"
)

func (h *Handler) BuildIndex() (err error) {
	fmt.Printf("Building index... ")
	start := time.Now()

	// TODO: build index
	h.index = &data.QuadNode{
		Centre: data.Point{
			Position: []float64{1, 2},
		},
		Children: []*data.QuadNode{{
			Centre: data.Point{
				Position: []float64{1, 2},
			},
		}, {
			Centre: data.Point{
				Position: []float64{1, 2},
			},
		}, {
			Centre: data.Point{
				Position: []float64{1, 2},
			},
		}, {
			Centre: data.Point{
				Position: []float64{1, 2},
			},
		}},
	}
	h.config.IndexReady = true
	lastExecTime := float64(time.Since(start).Nanoseconds())
	fmt.Printf("Done in %.3f µs (%.3f ms)\n", lastExecTime/1e3, lastExecTime/1e6)
	return
}

func countNodes(node *data.QuadNode) (count, depth int) {
	if node != nil {
		for _, child := range node.Children {
			c, d := countNodes(child)
			count += c
			if d > depth {
				depth = d
			}
		}
		count++
		depth++
	}
	return
}
