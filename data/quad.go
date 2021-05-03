package data

type QuadNode struct {
	Centre   Point // data if point
	PointID  *int  // nil if region
	Depth    int   // depth of node, only for region
	Children [8]*QuadNode
}
