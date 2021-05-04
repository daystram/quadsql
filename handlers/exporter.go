package handlers

import (
	"bufio"
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"

	"github.com/daystram/quadsql/data"
	"github.com/daystram/quadsql/db"
	"github.com/daystram/quadsql/utils"
)

const (
	SVG_PADDING = 16
	SVG_POINT_R = 4
)

func (h *Handler) DrawSVG(scale float64, filename string) {
	file, _ := os.Create(fmt.Sprintf("%s.svg", filename))
	defer file.Close()
	writer := bufio.NewWriter(file)
	s := svg.New(writer)
	s.Start(int((db.MAX_RANGE+2*SVG_PADDING)*scale), int((db.MAX_RANGE+2*SVG_PADDING)*scale))
	s.Scale(scale)
	s.Rect(0, 0, db.MAX_RANGE+2*SVG_PADDING, db.MAX_RANGE+2*SVG_PADDING, "fill:white;stroke:none")
	s.Rect(SVG_PADDING, SVG_PADDING, db.MAX_RANGE, db.MAX_RANGE, "fill:none;stroke:black;stroke-width:2")
	s.Text(convX(0), convY(-10), "(0,0)", "font-family:monospace;font-size:8")
	info := fmt.Sprintf("%d points", len(h.database.Table))
	if h.config.IsPointQuad {
		info += " @ Point Quad-tree"
	} else {
		info += " @ Region Quad-tree"
	}
	s.Text(convX(db.MAX_RANGE), convY(-10), info, "font-family:monospace;font-size:8;text-anchor:end")
	drawNode(s, h.config.IsPointQuad, h.index, [4]int{0, db.MAX_RANGE, db.MAX_RANGE, 0}, true) // bound: S, E, N, W
	s.Gend()
	s.End()
	writer.Flush()
}

func drawNode(s *svg.SVG, isPoint bool, node *data.QuadNode, bound [4]int, root bool) {
	if node != nil {
		x, y := int(node.Centre.Position[0]), int(node.Centre.Position[1])
		if node.PointID != nil {
			s.Circle(convX(x), convY(y), SVG_POINT_R, "fill:none;stroke:red;stroke-width:2")
		}
		if root {
			s.Polygon([]int{convX(x), convX(x + 4), convX(x - 4)}, []int{convY(bound[0]), convY(bound[0] - 6), convY(bound[0] - 6)}, "fill:blue;stroke:none")
			s.Polygon([]int{convX(bound[1]), convX(bound[1] + 6), convX(bound[1] + 6)}, []int{convY(y), convY(y + 4), convY(y - 4)}, "fill:blue;stroke:none")
			s.Polygon([]int{convX(x), convX(x + 4), convX(x - 4)}, []int{convY(bound[2]), convY(bound[2] + 6), convY(bound[2] + 6)}, "fill:blue;stroke:none")
			s.Polygon([]int{convX(bound[3]), convX(bound[3] - 6), convX(bound[3] - 6)}, []int{convY(y), convY(y + 4), convY(y - 4)}, "fill:blue;stroke:none")
		}
		hasChild := false
		for quad, child := range node.Children {
			hasChild = hasChild || child != nil
			var cBound [4]int
			switch quad {
			case 0: // NE
				cBound = [4]int{
					utils.Max(bound[0], y),
					bound[1],
					bound[2],
					utils.Max(bound[3], x),
				}
			case 2: // NW
				cBound = [4]int{
					utils.Max(bound[0], y),
					utils.Min(bound[1], x),
					bound[2],
					bound[3],
				}
			case 1: // SE
				cBound = [4]int{
					bound[0],
					bound[1],
					utils.Min(bound[2], y),
					utils.Max(bound[3], x),
				}
			case 3: // SW
				cBound = [4]int{
					bound[0],
					utils.Min(bound[1], x),
					utils.Min(bound[2], y),
					bound[3],
				}
			}
			drawNode(s, isPoint, child, cBound, false)
		}
		if hasChild {
			s.Line(convX(x), convY(bound[0]), convX(x), convY(bound[2]), "fill:none;stroke:blue")
			s.Line(convX(bound[1]), convY(y), convX(bound[3]), convY(y), "fill:none;stroke:blue")
		}
	}
}

func convX(x int) int {
	return x + SVG_PADDING
}

func convY(y int) int {
	return db.MAX_RANGE - y + SVG_PADDING
}
