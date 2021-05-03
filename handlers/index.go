package handlers

import (
	"fmt"
	"math"
	"time"

	"github.com/daystram/quadsql/data"
	"github.com/daystram/quadsql/db"
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
					// continue diving in
					node = &(*node).Children[getQuadrant((*node).Centre, point)]
				} else {
					// insert new node
					*node = &data.QuadNode{
						Centre:  point,
						PointID: id,
					}
					break
				}
			}
		}
	} else {
		// build Region index
		fmt.Printf("Building Region index... ")
		for i, point := range h.database.Table {
			id := new(int)
			*id = i
			if h.index == nil {
				depth := 1
				h.index = &data.QuadNode{
					Centre: data.Point{
						Position: []float64{db.MAX_RANGE / math.Pow(2, float64(depth)), db.MAX_RANGE / math.Pow(2, float64(depth))},
					},
					Depth: depth,
				}
				quad := getQuadrant(h.index.Centre, point)
				h.index.Children[quad] = &data.QuadNode{
					Centre:  point,
					PointID: id,
				}
			} else {
				node := h.index
				for {
					if node.Depth > 10 {
						break
					}
					quad := getQuadrant(node.Centre, point)
					if child := node.Children[quad]; child != nil {
						if child.PointID == nil {
							node = child
						} else {
							// create new internal node (subdiv)
							subdiv := createSubDiv(node.Centre, quad, node.Depth+1)
							// reinsert this leaf node (child)
							subdiv.Children[getQuadrant(subdiv.Centre, child.Centre)] = child
							node.Children[quad] = &subdiv
							// continue diving in
							node = node.Children[quad]
						}
					} else {
						// insert new leaf node
						node.Children[quad] = &data.QuadNode{
							Centre:  point,
							PointID: id,
						}
						break
					}
				}
			}
		}

	}

	lastExecTime := float64(time.Since(start).Nanoseconds())
	fmt.Printf("Done in %.3f µs (%.3f ms)\n", lastExecTime/1e3, lastExecTime/1e6)
	h.config.IndexReady = true
	return
}

func createSubDiv(parent data.Point, quad uint, depth int) data.QuadNode {
	/*
		  2 | 0
		 --- ---
		  3 | 1

		 10 | 00  xy
		---- ----
		 11 | 01

		 -+ | ++  xy
		---- ----
		 -- | +-
	*/
	centre := make([]float64, len(parent.Position))
	delta := db.MAX_RANGE / math.Pow(2, float64(depth))
	for dim, value := range parent.Position {
		// check each bits value for dimension comparison
		// must reverse due to lower dimension is at MSD (e.g. 0b101 -> xyz)
		if ((quad >> (len(centre) - 1 - dim)) & 1) == 1 {
			centre[dim] = value - delta // +ve bit -> less than parent's
		} else {
			centre[dim] = value + delta // +ve bit -> larger than parent's
		}
	}

	return data.QuadNode{
		Centre: data.Point{
			Position: centre,
		},
		Depth: depth,
	}
}

func getQuadrant(center, point data.Point) (quad uint) {
	/*
	  2 | 0
	 --- ---
	  3 | 1
	*/
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
