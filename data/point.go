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

func (p Point) Equals(other Point) bool {
	if len(p.Position) != len(other.Position) {
		return false
	}
	for i, coord := range p.Position {
		if math.Abs(coord-other.Position[i]) > epsilon {
			return false
		}
	}
	return true
}
