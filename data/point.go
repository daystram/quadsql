package data

import (
	"encoding/json"
	"strings"
)

type Point struct {
	Coordinate []float32
}

func ParsePoint(str string) (point Point, err error) {
	str = "[" + strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(str), "point("), ")") + "]"
	err = json.Unmarshal([]byte(str), &point.Coordinate)
	return
}

func (p Point) String() string {
	byteStr, _ := json.Marshal(&p.Coordinate)
	return "Point(" + strings.TrimSuffix(strings.TrimPrefix(string(byteStr), "["), "]") + ")"
}
