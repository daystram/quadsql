package handlers

import (
	"fmt"
	"time"

	"github.com/daystram/quadsql/data"
)

// h.index = &data.QuadNode{
// 	Centre: data.Point{
// 		Position: []float64{1, 2},
// 	},
// 	Children: []*data.QuadNode{{
// 		Centre: data.Point{
// 			Position: []float64{1, 2},
// 		},
// 	}, {
// 		Centre: data.Point{
// 			Position: []float64{1, 2},
// 		},
// 	}, {
// 		Centre: data.Point{
// 			Position: []float64{1, 2},
// 		},
// 	}, {
// 		Centre: data.Point{
// 			Position: []float64{1, 2},
// 		},
// 	}},
// }

func (h *Handler) BuildIndex(isPoint bool) (err error) {
	start := time.Now()
	h.index = nil
	h.config.IsPointQuad = isPoint
	if h.config.IsPointQuad {
		// build Point index
		fmt.Printf("Building Point index... ")
		for i, point := range h.database.Table {
			id := new(int)
			*id = i
			node := &h.index
			for {
				if *node != nil {
					node = &(*node).Children[getQuadrant((*node).Centre, point)]
				} else {
					*node = &data.QuadNode{
						Centre:  point,
						PointID: id,
					}
					break
				}
			}
		}
	} else {
		// TODO: build Region index
		fmt.Printf("Building Region index... ")

	}

	lastExecTime := float64(time.Since(start).Nanoseconds())
	fmt.Printf("Done in %.3f µs (%.3f ms)\n", lastExecTime/1e3, lastExecTime/1e6)
	h.config.IndexReady = true
	return
}

func getQuadrant(center, point data.Point) (quad uint) {
	for i := range center.Position {
		quad = quad << 1
		if center.Position[i] > point.Position[i] {
			quad |= 1
		}
	}
	return quad
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
