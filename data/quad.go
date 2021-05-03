package data

type QuadNode struct {
	Centre   Point // point: internal, leave; region: internal
	PointID  *int  // point: internal, leave; region: leave
	Children [8]*QuadNode
}
