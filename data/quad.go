package data

type QuadNode struct {
	Centre   Point // data if point
	PointID  *int  // nil if region
	Children [8]*QuadNode
}
