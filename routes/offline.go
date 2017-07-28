package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	config "github.com/traderboy/collector-server/config"
)

func offline(w http.ResponseWriter, r *http.Request) {
	dbPath := r.URL.Query().Get("db")
	if len(dbPath) == 0 {
		log.Println("No database entered")
		return
	}
	log.Println("/offline/ (" + r.Method + ")")
	log.Println("Database: " + dbPath)
	dbName := "file:" + dbPath + config.SqlWalFlags //"?PRAGMA journal_mode=WAL"

	DbCollectorDb, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	sqlstr := "SELECT 'json','GDB_ServiceItems','ItemInfo','DatasetName',\"DatasetName\",\"ItemType\",\"ItemId\" from \"GDB_ServiceItems\" UNION SELECT 'xml','GDB_Items','Definition','Name',\"Name\",\"ObjectID\",\"DatasetSubtype1\" FROM " + config.Schema + config.DblQuote("GDB_Items")
	log.Printf("Query: " + sqlstr)

	/*
		stmt, err := DbCollectorDb.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
		}
	*/
	//rows := stmt.QueryRow(id)

	rows, err := DbCollectorDb.Query(sqlstr)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		w.Write(response)
		return
	}
	//defer rows.Close()

	var table []byte
	var format []byte
	var value []byte
	var field []byte
	var queryField []byte
	var itemtype []byte
	var itemid []byte

	var results [][]string
	//var items map[string]interface{}

	for rows.Next() {
		err := rows.Scan(&table, &format, &field, &queryField, &value, &itemtype, &itemid)
		if err != nil {
			log.Fatal(err)
		}
		vals := []string{string(table), string(format), string(field), string(queryField), string(value), string(itemtype), string(itemid)}
		results = append(results, vals)

	}
	rows.Close()
	DbCollectorDb.Close()

	//err = DbCollectorDb.QueryRow(sql).Scan(&itemInfo)
	//rows, err := Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)

	//DbCollectorDb.Close()
	response, _ := json.Marshal(results)

	//response, _ := json.Marshal(map[string]interface{}itemInfo)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func offline_load(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//value := vars["value"]
	//dbPath := r.URL.Query().Get("db")
	//vars := mux.Vars(r)
	vars := mux.Vars(r)
	table := vars["table"]
	field := vars["field"]
	queryField := vars["queryField"]
	value := vars["value"]
	if value == " " {
		value = ""
	}
	format := vars["format"]

	//id := vars["id"]
	//idInt, _ := strconv.Atoi(id)
	dbPath := r.URL.Query().Get("db")
	if len(dbPath) == 0 {
		log.Println("No database entered")
		return
	}
	log.Println("/offline/" + table + "/" + field + "/" + queryField + "/" + value + " (" + r.Method + ")")
	//tableName := config.Collector.Projects[name].Layers[id]["data"].(string)
	//tableName = strings.ToUpper(tableName)
	//log.Println("/offline/"+type+"/"+name)
	dbName := "file:" + dbPath + config.SqlWalFlags //+ "?PRAGMA journal_mode=WAL"
	DbCollectorDb, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	if r.Method == "PUT" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
			return
		}

		//ret := config.SetArcService(body, name, "FeatureServer", idInt, "")
		sql := "update " + config.DblQuote(table) + " set " + config.DblQuote(field) + "=? where " + queryField + "=?" //OBJECTID=?"
		stmt, err := DbCollectorDb.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)

			return
		}
		log.Println("Updating table: " + value)
		log.Println(sql)
		//log.Println(strings.Replace(string(body), "'", "''", -1))

		_, err = stmt.Exec(strings.Replace(string(body), "'", "''", -1), value)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
			return
		}

		//db.Close()
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": "ok"})
		w.Write(response)
		return
	}
	if format == "table" {
		sql := "SELECT * FROM " + config.DblQuote(value)
		log.Printf("Query: " + sql)

		//var itemInfo *[]byte
		//*interface{}
		rows, err := DbCollectorDb.Query(sql)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		// get the column names from the query
		var columns []string
		columns, err = rows.Columns()
		colNum := len(columns)
		//<style>table{width:100%;}table, th, td { border: 1px solid black;  border-collapse: collapse;}th, td { padding: 5px; text-align: left;}</style>
		t := "<table class='table-bordered table-striped'>"
		for n := 0; n < colNum; n++ {
			t = t + "<th>" + columns[n] + "</th>"
		}
		rawResult := make([][]byte, colNum)
		for rows.Next() {
			cols := make([]interface{}, colNum)
			for i := 0; i < colNum; i++ {
				cols[i] = &rawResult[i]
			}
			err = rows.Scan(cols...)
			if err != nil {
				log.Println(err.Error())
				//w.Header().Set("Content-Type", "application/json")
				//response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
				//w.Write(response)
				w.Write([]byte("Error: " + err.Error()))

				return
			}
			t = t + "<tr>"
			//for i := 0; i < colNum; i++ {
			for i, raw := range rawResult {
				//w.Write(cols[i])
				if strings.ToLower(columns[i]) == "shape" {
					t = t + "<td>Shape</td>"
				} else {
					t = t + fmt.Sprintf("<td>%v</td>", string(raw))
				}
				//w.Write([]byte(cols[i]))
			}
			t = t + "</tr>"
			//s := fmt.Sprintf("a %s", "string")
			//w.Write([]byte(s))
			//for i := 0; i < colNum; i++ {
			//cols[i] = VehicleCol(columns[i], &vh)
			//w.Write(rows.Scan(cols...)
			//}
			//err = rows.Scan(&itemInfo)

			//for num, i := range *itemInfo {
			//	w.Write(i)
			//}
		}
		t = t + "</table>"
		w.Write([]byte(t))
		return
	}
	//Db.Exec(initializeStr)
	//log.Print("Sqlite database: " + dbName)
	//sql := "SELECT \"DatasetName\",\"ItemId\",\"ItemInfo\",\"AdvancedDrawingInfo\" FROM \"GDB_ServiceItems\""
	sql := "SELECT " + config.DblQuote(field) + " from " + config.DblQuote(table) + " where " + config.DblQuote(queryField) + "=" + config.DblQuote(value) // 'GDB_ServiceItems',\"DatasetName\" from \"GDB_ServiceItems\" UNION SELECT 'GDB_Items',\"Name\" FROM \"GDB_Items\""
	log.Printf("Query: " + sql)
	stmt, err := DbCollectorDb.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		w.Write(response)
	}
	//rows := stmt.QueryRow(id)
	var itemInfo []byte
	err = stmt.QueryRow().Scan(&itemInfo)
	//rows, err := Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		w.Write(response)
		return
	}
	stmt.Close()
	DbCollectorDb.Close()
	//response, _ := json.Marshal(map[string]interface{}itemInfo)
	if format == "xml" {
		w.Header().Set("Content-Type", "application/xml")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	w.Write(itemInfo)
	//w.Header().Set("Content-Type", "application/xml")
	//w.Write(itemInfo)
	/*
		var id int
		idstr := r.URL.Query().Get("id")
		dbPath := r.URL.Query().Get("db")

		if len(idstr) > 0 {
			id, _ = strconv.Atoi(idstr)
		} else {
			id = config.DbSource
		}
	*/
}
