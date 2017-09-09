package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
	"os"
	"strconv"
)

func unsigned_to_signed(i int64, max_positive int64) (int64) {
	if i < max_positive {
		return i
	}
	return (i - (max_positive * 2))
}

func pythonmodulo(i int64, mod int64) (int64) {
	if i >= 0 {
		return i % mod
	}
	return (mod - ((-i) % mod))
}

func main() {
	var ci, co int

	if len(os.Args) != 3 {
		log.Fatal("Not enough arguments: <sqlite file> <cutoff limit>")
	}
	f := os.Args[1]
	l, err := strconv.ParseInt(os.Args[2], 0, 64)
	if err != nil {
		log.Fatal(err)
	}
	if l < 0 {
		log.Fatal("cutoff limit should be positive")
	}

	db, err := sql.Open("sqlite3", f)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select pos from blocks")
	if err != nil {
		log.Fatal(err)
	}

	var arr []int64

	for rows.Next() {
		var pos int64
		err = rows.Scan(&pos)
		if err != nil {
			log.Fatal(err)
		}
		arr = append(arr, int64(pos))
		ci++
	}
	rows.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("delete from blocks where pos = ?")
	if err != nil {
		log.Fatal(err)
	}

	for _, pos := range arr {
		var opos = pos

		var x = unsigned_to_signed(pythonmodulo(pos, 4096), 2048)
		pos = (pos - x) / 4096
		var y = unsigned_to_signed(pythonmodulo(pos, 4096), 2048)
		pos = (pos - y) / 4096
		var z = unsigned_to_signed(pythonmodulo(pos, 4096), 2048)
		if (x * 16 > l) || (x * 16 < -l) || (z * 16 > l) || (z * 16 < -l) {
			_, err = stmt.Exec(fmt.Sprintf("%v", opos))
			if err != nil {
				log.Fatal(err)
			}
			co++
		}
	}
	tx.Commit()

	_, err = db.Exec("VACUUM")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("map.sqlite: removed %v of %v blocks\n", co, ci)

	defer db.Close()
}

