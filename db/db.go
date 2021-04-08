package db

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/daystram/quadsql/data"
)

type DB struct {
	Dimension int
	Table     []data.Point // in-memory table
	file      *os.File
}

func OpenDB(source string) (db DB, err error) {
	if db.file, err = os.Open(source); err != nil {
		return
	}
	// parse source
	scanner := bufio.NewScanner(db.file)
	scanner.Scan()
	dimensionInfo := scanner.Text()
	if db.Dimension, err = strconv.Atoi(strings.Split(dimensionInfo, " ")[1]); err != nil {
		fmt.Printf("invalid dimension definition: \"%s\"\n", dimensionInfo)
		return DB{}, ErrBadDBSource
	}
	db.Table = make([]data.Point, 0)
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
		db.Table = append(db.Table, point)
	}
	return
}

func (db *DB) Close() error {
	return db.file.Close()
}
