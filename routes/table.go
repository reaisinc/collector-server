package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	config "github.com/traderboy/collector-server/config"
	structs "github.com/traderboy/collector-server/structs"
)

func table_id(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	var tableName string
	_, err1 := strconv.Atoi(id)
	if err1 != nil {
		tableName = id

	} else {
		tableName = config.Collector.Projects[name].Layers[id].Data
	}
	dbPath := r.URL.Query().Get("db")

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/table/" + id)
	var dbName = "file:" + config.Collector.Projects[name].ReplicaPath // + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase" + config.SqlWalFlags
	if len(dbPath) > 0 {
		if config.Collector.DefaultDataSource != structs.PGSQL {
			if config.DbSqliteDbName != dbPath {
				if config.DbSqliteQuery != nil {
					config.DbSqliteQuery.Close()
				}
				config.DbSqliteQuery = nil
			}
			config.DbSqliteDbName = dbPath
			dbName = "file:" + dbPath + config.SqlWalFlags //+ "?PRAGMA journal_mode=WAL"
		}
		/*
			if config.DbSqliteQuery != nil {
				config.DbSqliteQuery.Close()
				config.DbSqliteQuery = nil
				if dbPath == "close" {
					return
				}
			}
		*/
	} else {
		if config.DbSqliteDbName != dbName {
			if config.DbSqliteQuery != nil {
				config.DbSqliteQuery.Close()
			}
			config.DbSqliteQuery = nil
		}
		config.DbSqliteDbName = dbName
	}

	var err error
	//if err != nil {
	if config.DbSqliteQuery == nil {
		//config.DbSqliteQuery, err = sql.Open("sqlite3", "file:"+dbName+"?PRAGMA journal_mode=WAL")
		if config.Collector.DefaultDataSource == structs.PGSQL {
			config.DbSqliteQuery = config.GetReplicaDB(name)
		} else {
			config.DbSqliteQuery, err = sql.Open("sqlite3", dbName)
			if err != nil {
				log.Fatal(err)
				w.Write([]byte("Error: " + err.Error()))
				return
			}
		}
	}

	//Db.Exec(initializeStr)
	log.Print("Sqlite database: " + dbName)
	//sql := "SELECT \"DatasetName\",\"ItemId\",\"ItemInfo\",\"AdvancedDrawingInfo\" FROM \"GDB_ServiceItems\""
	sql := "SELECT * FROM " + config.DblQuote(tableName)
	log.Printf("Query: " + sql)

	//var itemInfo *[]byte
	//*interface{}
	rows, err := config.DbSqliteQuery.Query(sql)
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
	//.Scan(&itemInfo)
	//rows, err := Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)

	//w.Header().Set("Content-Type", "application/xml")
	//w.Write(itemInfo)
}
