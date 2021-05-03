package handlers

import (
	"bufio"
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"

	"github.com/daystram/quadsql/data"
	"github.com/daystram/quadsql/db"
)

const (
	SVG_PADDING = 16
	SVG_POINT_R = 4
)

func DrawSVG(isPoint bool, node *data.QuadNode, scale float64, filename string) {
	file, _ := os.Create(fmt.Sprintf("%s.svg", filename))
	defer file.Close()
	writer := bufio.NewWriter(file)
	s := svg.New(writer)
	s.Start(db.MAX_RANGE+2*SVG_PADDING, db.MAX_RANGE+2*SVG_PADDING)
	s.Scale(scale)
	s.Rect(SVG_PADDING, SVG_PADDING, db.MAX_RANGE, db.MAX_RANGE, "fill:none;stroke:black;stroke-width:2")
	drawNode(s, isPoint, node, [4]int{SVG_PADDING, db.MAX_RANGE + SVG_PADDING, db.MAX_RANGE + SVG_PADDING, SVG_PADDING}, true)
	s.Gend()
	s.End()
	writer.Flush()
}

func drawNode(s *svg.SVG, isPoint bool, node *data.QuadNode, bound [4]int, root bool) {
	if node != nil {
		x, y := int(node.Centre.Position[0])+SVG_PADDING, int(node.Centre.Position[1])+SVG_PADDING
		if node.PointID != nil {
			s.Circle(x, y, SVG_POINT_R, "fill:none;stroke:red;stroke-width:2")
		}
		if root {
			s.Polygon([]int{x, x + 4, x - 4}, []int{bound[0], bound[0] - 6, bound[0] - 6}, "fill:blue;stroke:none")
			s.Polygon([]int{bound[1], bound[1] + 6, bound[1] + 6}, []int{y, y + 4, y - 4}, "fill:blue;stroke:none")
			s.Polygon([]int{x, x + 4, x - 4}, []int{bound[2], bound[2] + 6, bound[2] + 6}, "fill:blue;stroke:none")
			s.Polygon([]int{bound[3], bound[3] - 6, bound[3] - 6}, []int{y, y + 4, y - 4}, "fill:blue;stroke:none")
		}
		hasChild := false
		for quad, child := range node.Children {
			hasChild = hasChild || child != nil
			/*
				Mirrorred quadrants (origin on top left):
				 3 | 1
				--- ---
				 2 | 0
			*/
			switch quad {
			case 1: // NE
				drawNode(s, isPoint, child, [4]int{
					bound[0],
					bound[1],
					min(bound[2], y),
					max(bound[3], x),
				}, false)
			case 3: // NW
				drawNode(s, isPoint, child, [4]int{
					bound[0],
					min(bound[1], x),
					min(bound[2], y),
					bound[3],
				}, false)
			case 0: // SE
				drawNode(s, isPoint, child, [4]int{
					max(bound[0], y),
					bound[1],
					bound[2],
					max(bound[3], x),
				}, false)
			case 2: // SW
				drawNode(s, isPoint, child, [4]int{
					max(bound[0], y),
					min(bound[1], x),
					bound[2],
					bound[3],
				}, false)
			}
		}
		if hasChild {
			s.Line(x, bound[0], x, bound[2], "fill:none;stroke:blue")
			s.Line(bound[3], y, bound[1], y, "fill:none;stroke:blue")
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
