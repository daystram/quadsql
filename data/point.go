package data

import (
	"encoding/json"
	"math"
	"strings"
)

const epsilon = 1e-9

type Point struct {
	Position []float64
}

func ParsePoint(str string) (point Point, err error) {
	str = strings.TrimPrefix(strings.ToLower(str), "point")
	str = "[" + strings.TrimSuffix(strings.TrimPrefix(str, "("), ")") + "]"
	err = json.Unmarshal([]byte(str), &point.Position)
	return
}

func (p Point) String() string {
	byteStr, _ := json.Marshal(&p.Position)
	return "Point(" + strings.TrimSuffix(strings.TrimPrefix(string(byteStr), "["), "]") + ")"
}

func (p Point) CompareTo(other Point) (diff float64) {
	if len(p.Position) > len(other.Position) {
		return -1
	}
	if len(p.Position) < len(other.Position) {
		return 1
	}
	inMargin := true
	var delta float64
	for dim := 0; dim < len(p.Position); dim++ {
		delta += other.Position[dim] - p.Position[dim]
		inMargin = inMargin && math.Abs(delta) < epsilon
	}
	if inMargin {
		return 0
	} else {
		return delta
	}
}
