package data

import (
	"encoding/json"
	"strings"
)

type Point struct {
	Coordinate []float32
}

func ParsePoint(strPoint string) (point Point, err error) {
	strPoint = "[" + strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(strPoint), "point("), ")") + "]"
	err = json.Unmarshal([]byte(strPoint), &point.Coordinate)
	return
}

func (p Point) String() string {
	byteStr, _ := json.Marshal(&p.Coordinate)
	return "Point(" + strings.TrimSuffix(strings.TrimPrefix(string(byteStr), "["), "]") + ")"
}
