//
//    mtmapprune - Prune blocks outside a limit box in map.sqlite
//
//    Copyright (C) 2017 - Auke Kok <sofar@foo-projects.org>
//    Portions Copyright (C) 2013 celeron55, Perttu Ahola <celeron55@gmail.com>
//
//    This library is free software; you can redistribute it and/or
//    modify it under the terms of the GNU Lesser General Public
//    License as published by the Free Software Foundation; either
//    version 2.1 of the License, or (at your option) any later version.
//
//    This library is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
//    Lesser General Public License for more details.
//
//    You should have received a copy of the GNU Lesser General Public
//    License along with this library; if not, write to the Free Software
//    Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
//

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // MIT licensed.
	"log"
	"os"
	"strconv"
)

// From: minetest/src/database-sqlite3.cpp
func unsigned_to_signed(i int64, max_positive int64) int64 {
	if i < max_positive {
		return i
	}
	return (i - (max_positive * 2))
}

func parse_arg(n int64) int64 {
	v, err := strconv.ParseInt(os.Args[n], 0, 64)
	if err != nil {
		log.Fatal("Error parsing argument: ", n, ": ", err)
	}
	return v
}

type Limit struct {
	min int64
	max int64
}

func main() {
	var ci, co int

	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments: sqlite_file max_x [max_y [max_z [min_x min_y min_z]]]")
	}
	if len(os.Args) == 6 || len(os.Args) == 7 {
		log.Fatal("Argument count must be 4 or 7, if you specify one min limit, you must specify them all")
	}

	f := os.Args[1]

	var x Limit
	var y Limit
	var z Limit

	x.max = parse_arg(2)
	if len(os.Args) > 3 {
		y.max = parse_arg(3)
		if len(os.Args) > 4 {
			z.max = parse_arg(4)
		} else {
			z.max = x.max
		}
	} else {
		y.max = x.max
		z.max = x.max
	}

	if len(os.Args) == 8 {
		x.min = parse_arg(5)
		y.min = parse_arg(6)
		z.min = parse_arg(7)
	} else {
		if (x.max < 0) || (y.max < 0) || (z.max < 0) {
			log.Fatal("Limits should be positive when passing max_x max_y max_z values")
		}
		x.min = -x.max
		y.min = -y.max
		z.min = -z.max
	}

	// ensure proper ordering of min, max
	if x.min > x.max {
		a := x.min
		x.min = x.max
		x.max = a
	}
	if y.min > y.max {
		a := y.min
		y.min = y.max
		y.max = a
	}
	if z.min > z.max {
		a := z.min
		z.min = z.max
		z.max = a
	}

	db, err := sql.Open("sqlite3", f)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Collecting blocks")
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

	log.Println("Deleting blocks")
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

		// From: minetest/src/database-sqlite3.cpp
		var xi = unsigned_to_signed(pos&0xfff, 2048)
		pos = (pos - xi) / 4096
		var yi = unsigned_to_signed(pos&0xfff, 2048)
		pos = (pos - yi) / 4096
		var zi = unsigned_to_signed(pos&0xfff, 2048)
		if (xi*16 > x.max) || (xi*16+15 < x.min) ||
			(yi*16 > y.max) || (yi*16+15 < y.min) ||
			(zi*16 > z.max) || (zi*16+15 < z.min) {
			_, err = stmt.Exec(fmt.Sprintf("%v", opos))
			if err != nil {
				log.Fatal(err)
			}
			co++
		}
	}
	tx.Commit()

	log.Println("Vaccuuming database")
	_, err = db.Exec("VACUUM")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Removed %v of %v blocks (limits: [%v, %v, %v]-[%v, %v, %v])\n",
		co, ci, x.min, y.min, z.min, x.max, y.max, z.max)

	defer db.Close()
}
