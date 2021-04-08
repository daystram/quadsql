package data

type QuadNode struct {
	Centre   Point // point: internal, leave; region: internal
	PointID  int64 // point: internal, leave; region: leave
	Children []*QuadNode
}
