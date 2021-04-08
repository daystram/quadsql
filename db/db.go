package db

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type DB struct {
	Dimension int
	RowCount  int
	file      *os.File
}

func OpenDB(source string) (db DB, err error) {
	if db.file, err = os.Open(source); err != nil {
		return
	}
	// validate source
	scanner := bufio.NewScanner(db.file)
	scanner.Scan()
	dimensionInfo := scanner.Text()
	if db.Dimension, err = strconv.Atoi(strings.Split(dimensionInfo, " ")[1]); err != nil {
		fmt.Printf("invalid dimension definition: \"%s\"\n", dimensionInfo)
		return DB{}, ErrBadDBSource
	}
	r, _ := regexp.Compile("^point( ([0-9]+[.])?[0-9]+){2,}$")
	for scanner.Scan() {
		line := scanner.Text()
		if !r.MatchString(line) {
			fmt.Printf("invalid statement: \"%s\"\n", line)
			return DB{}, ErrBadDBSource
		}
	}
	return
}

func (db *DB) Close() error {
	return db.file.Close()
}
