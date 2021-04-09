package db

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/daystram/quadsql/data"
)

const MAX_RANGE = 1024

type DB struct {
	Dimension int
	Table     map[int]data.Point // in-memory table
	file      *os.File
}

func OpenDB(source string) (db DB, err error) {
	if db.file, err = os.OpenFile(source, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return
	}
	// parse source
	scanner := bufio.NewScanner(db.file)
	scanner.Scan()
	dimensionInfo := scanner.Text()
	if dimensionInfo != "" {
		if db.Dimension, err = strconv.Atoi(strings.Split(dimensionInfo, " ")[1]); err != nil {
			fmt.Printf("invalid dimension definition: \"%s\"\n", dimensionInfo)
			return DB{}, ErrBadDBSource
		}
	}
	// populate table
	lastID := 0
	db.Table = make(map[int]data.Point)
	for scanner.Scan() {
		line := scanner.Text()
		var point data.Point
		if point, err = data.ParsePoint(line); err != nil {
			fmt.Printf("invalid statement: \"%s\"\n", line)
			return DB{}, ErrBadDBSource
		}
		if len(point.Position) != db.Dimension {
			fmt.Printf("point %s does not match DB dimension of %d\n", point, db.Dimension)
			return DB{}, ErrBadDBSource
		}
		db.Table[lastID] = point
		lastID++
	}
	return
}

func (db *DB) Close() error {
	db.file.Truncate(0)
	db.file.Seek(0, 0)
	writer := bufio.NewWriter(db.file)
	fmt.Fprintf(writer, "dim %d\n", db.Dimension)
	points := make([]data.Point, len(db.Table))
	for id, point := range db.Table {
		points[id] = point
	}
	for _, point := range points {
		fmt.Fprintf(writer, "%s\n", point)
	}
	writer.Flush()
	return db.file.Close()
}
