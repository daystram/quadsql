package handlers

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/daystram/quadsql/data"
	"github.com/daystram/quadsql/db"
)

func Generate(source string, genSeed int64, dimension, size int, max float64, sorted bool) (err error) {
	var database db.DB
	if database, err = db.OpenDB(source); err != nil {
		return
	}

	start := time.Now()
	database.Dimension = dimension
	randomizer := rand.New(rand.NewSource(genSeed))
	points := make([]data.Point, 0)
	for id := 0; id < size; id++ {
		position := make([]float64, 0)
		for c := 0; c < dimension; c++ {
			position = append(position, randomizer.Float64()*max)
		}
		points = append(points, data.Point{Position: position})
	}
	if sorted {
		sort.Slice(points, func(i, j int) bool {
			return points[i].CompareTo(points[j]) > 0
		})
	}
	database.Table = make(map[int]data.Point)
	for id, point := range points {
		database.Table[id] = point
	}
	execTime := float64(time.Since(start).Nanoseconds())

	err = database.Close()
	fmt.Printf("Exec time: %.3f ms\n", execTime/1e6)
	return
}
