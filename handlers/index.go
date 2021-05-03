package handlers

import (
	"fmt"
	"time"

	"github.com/daystram/quadsql/data"
	"github.com/daystram/quadsql/db"
)

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
						Centre:   point,
						PointID:  id,
						Children: make([]*data.QuadNode, Exp2(h.database.Dimension)),
					}
					break
				}
			}
		}
	} else {
		// build Region index
		fmt.Printf("Building Region index... ")
		h.index = &data.QuadNode{
			Centre: data.Point{
				Position: []float64{db.MAX_RANGE / 2, db.MAX_RANGE / 2},
			},
			Children: make([]*data.QuadNode, Exp2(h.database.Dimension)),
		}
		for i, point := range h.database.Table {
			id := new(int)
			*id = i
			node := &h.index
			// track parent and depth for subdivision
			parent, depth := h.index, 0
			for {
				if *node != nil {
					if (*node).PointID == nil {
						// continue diving in
						depth++
						parent = *node
						node = &(*node).Children[getQuadrant((*node).Centre, point)]
					} else {
						// collision with leaf node: create subdivision and reinsert leaf node
						subdiv := createSubDiv(parent.Centre, getQuadrant(parent.Centre, point), depth+1)
						subdiv.Children[getQuadrant(subdiv.Centre, (*node).Centre)] = *node
						*node = &subdiv
					}
				} else {
					// insert new leaf node
					*node = &data.QuadNode{
						Centre:   point,
						PointID:  id,
						Children: make([]*data.QuadNode, Exp2(h.database.Dimension)),
					}
					break
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
	dimension := len(parent.Position)
	centre := make([]float64, dimension)
	delta := db.MAX_RANGE / float64(Exp2(depth))
	for dim, value := range parent.Position {
		// check each bits' value for dimension comparison
		// must reverse due to lower dimension is at MSD (e.g. 0b101 -> xyz)
		if ((quad >> (dimension - 1 - dim)) & 1) == 1 {
			centre[dim] = value - delta // +ve bit -> less than parent's
		} else {
			centre[dim] = value + delta // +ve bit -> larger than parent's
		}
	}

	return data.QuadNode{
		Centre: data.Point{
			Position: centre,
		},
		Children: make([]*data.QuadNode, Exp2(dimension)),
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

func Exp2(x int) int {
	y := 1
	for i := 0; i < x; i++ {
		y *= 2
	}
	return y
}
