package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/handlers"
	config "github.com/traderboy/collector-server/config"
	routes "github.com/traderboy/collector-server/routes"
)

var logPath = "logfile.txt"

func main() {
	test()
	return
}

func test() {

	//config.DbSqliteQuery, err = sql.Open("sqlite3", "file:"+dbName+"?PRAGMA journal_mode=WAL")
	dbName := "C:\\docker\\src\\github.com\\traderboy\\collector-server\\catalogs\\leasecompliance2016\\replicas\\leasecompliance2016.geodatabase"
	Db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	sql := "select shape from homesites"
	var fields []byte
	err = Db.QueryRow(sql).Scan(&fields)
	f := fields[len(fields)-9:]
	fmt.Println("%v", fields)
	fmt.Println("%v", f)
	fmt.Println("%v", Float64frombytes(f))
	for index, _ := range fields {
		//fmt.Println("%v: %v", index, b)
		fmt.Println("%v: %v:  %v", index, fields[index:index+8], Float64frombytes(fields[index:index+8]))

	}

	return

	var pi float64
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	b = fields[1:4]
	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &pi)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	//fmt.Print(pi)

	//var v uint32
	//err := binary.Read(bytes.NewReader(fields), &v)

	//data := binary.BigEndian.Uint64(fields[:4])
	//fmt.Println(data)
}
func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
