package data

import (
	"encoding/json"
	"strings"
)

type Point struct {
	Position []float32
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
