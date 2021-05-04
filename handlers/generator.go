package handlers

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/daystram/quadsql/data"
	"github.com/daystram/quadsql/db"
)

func Generate(source string, genSeed int64, dimension, size int, max float64, distribution string, sorted bool) (err error) {
	var database db.DB
	if database, err = db.OpenDB(source); err != nil {
		return
	}

	start := time.Now()
	database.Dimension = dimension
	randomizer := rand.New(rand.NewSource(genSeed))
	database.Table = make([]data.Point, 0)
	id := 0
	for id < size {
		var position []float64
		switch distribution {
		case "uniform":
			position = uniform(randomizer, dimension, max)
		case "normal":
			position = normal(randomizer, dimension, max, max/2, max/6)
		case "line":
			position = line(randomizer, dimension, max, max/16)
		case "line-strict":
			position = line(randomizer, dimension, max, 0)
		case "exp":
			position = exponential(randomizer, dimension, max, 2/max)
		default:
			return fmt.Errorf("unsupported. supported distributions: uniform, normal, line, line-strict, exp")
		}
		ok := true
		for _, value := range position {
			if value < 0 || value > max {
				ok = false
			}
		}
		if ok {
			database.Table = append(database.Table, data.Point{Position: position})
			id++
		}
	}
	if sorted {
		sort.Slice(database.Table, func(i, j int) bool {
			return database.Table[i].CompareTo(database.Table[j]) > 0
		})
	}
	execTime := float64(time.Since(start).Nanoseconds())

	err = database.Close()
	fmt.Printf("Exec time: %.3f ms\n", execTime/1e6)
	return
}

func uniform(randomizer *rand.Rand, dimension int, max float64) (position []float64) {
	for c := 0; c < dimension; c++ {
		position = append(position, randomizer.Float64()*max)
	}
	return
}

func normal(randomizer *rand.Rand, dimension int, max, mean, sd float64) (position []float64) {
	for c := 0; c < dimension; c++ {
		position = append(position, randomizer.NormFloat64()*sd+mean)
	}
	return
}

func line(randomizer *rand.Rand, dimension int, max, sd float64) (position []float64) {
	center := randomizer.Float64() * max
	for c := 0; c < dimension; c++ {
		position = append(position, center+randomizer.NormFloat64()*sd)
	}
	return
}

func exponential(randomizer *rand.Rand, dimension int, max, lambda float64) (position []float64) {
	for c := 0; c < dimension; c++ {
		position = append(position, randomizer.ExpFloat64()/lambda)
	}
	return
}
