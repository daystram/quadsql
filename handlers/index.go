package handlers

import (
	"fmt"

	"github.com/daystram/quadsql/data"
)

func (h *Handler) BuildIndex() (err error) {
	fmt.Printf("Building index... ")
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
	fmt.Println("Done")
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
