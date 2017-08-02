package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	//structs "github.com/traderboy/collector-server/structs"
	config "github.com/traderboy/collector-server/config"
)

func db(w http.ResponseWriter, r *http.Request) {
	log.Println("/db (" + r.Method + ")")
	//vars := mux.Vars(r)
	str := "<ul>"
	for _, val := range config.Collector.Projects {
		str += "<li>" + val.Name + "</li>"
		//fmt.Printf("%v: %v\n", key, val)
		//fmt.Println()
		/*
			for _, val1 := range val.Layers {
				//fmt.Printf("%v: %v\n", key1, val1)
				//fmt.Println()
				for key2, val2 := range val1 {
					log.Printf("%v: %v\n", key2, val2)
					//fmt.Println()
				}

			}
		*/
	}
	str += "</ul>"

	//log.Println(Collector.Projects["layers"].ReplicaPath)
	//log.Println(Collector.Projects["layers"].Layers["0"]["oidname"].(string))
	/*
		var id int
		idstr := r.URL.Query().Get("id")

		if len(idstr) > 0 {
			id, _ = strconv.Atoi(idstr)
		} else {
			id = config.DbSource
		}
	*/
	//	PGSQL   = 1
	//	SQLITE3 = 2
	//	FILE    = 3
	/*
		if id == 3 {
			str += "<li>Static JSON files <b style='color:red'>active </b></li>"
			config.SetDatasource(structs.FILE)
		} else {
			str += "<li>Static JSON files <a href='/db?id=3'>enable</a> </li>"
		}
		if id == 2 {
			str += "<li>Sqlite <b style='color:red'>active </b> </li>"
			config.SetDatasource(structs.SQLITE3)
		} else {
			str += "<li>Sqlite <a href='/db?id=2'>enable</a> </li>"
		}
		if id == 1 {
			str += "<li>Postgresql <b style='color:red'>active </b> </li>"
			config.SetDatasource(structs.PGSQL)
		} else {
			str += "<li>Postgresql <a href='/db?id=1'>enable</a> </li>"
		}
	*/
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Data source</h1><ul>" + str + "</ul>"))

}

func db_id(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	idInt, _ := strconv.Atoi(id)
	idInt++
	//newId := strconv.Itoa(idInt + 1)
	fieldStr := r.URL.Query().Get("field")
	if len(fieldStr) == 0 {
		fieldStr = config.DblQuote("ItemInfo")
	}
	dbPath := r.URL.Query().Get("db")
	if len(dbPath) > 0 {
	}

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/db/" + id)

	//var dbName = config.Collector.Projects[name].ReplicaPath // + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
	//var parentObjectID = config.Collector.Projects[name].Layers[id]["oidname"].(string)
	/*
		if len(dbPath) > 0 {
			if config.DbSqliteDbName != dbPath {
				if config.DbSqliteQuery != nil {
					config.DbSqliteQuery.Close()
				}
				config.DbSqliteQuery = nil
			}
			config.DbSqliteDbName = dbPath
			dbName = "file:" + dbPath + config.SqlWalFlags //"?PRAGMA journal_mode=WAL"
		} else {
			if config.DbSqliteDbName != dbName {
				if config.DbSqliteQuery != nil {
					config.DbSqliteQuery.Close()
				}
				config.DbSqliteQuery = nil
			}
			config.DbSqliteDbName = dbName
		}
	*/
	//err := config.DbSqliteQuery.Ping()

	var err error
	//if err != nil {
	/*
		if config.DbSqliteQuery == nil {
			//config.DbSqliteQuery, err = sql.Open("sqlite3", "file:"+dbName+"?PRAGMA journal_mode=WAL")
			config.DbSqliteQuery, err = sql.Open("sqlite3", dbName)
			if err != nil {
				log.Fatal(err)
			}
		}
	*/
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
			return
		}
		//ret := config.SetArcService(body, name, "FeatureServer", idInt, "")
		sql := "update " + config.Collector.Schema + config.DblQuote("GDB_ServiceItems") + " set " + fieldStr + "=? where " + config.DblQuote("OBJECTID") + "=?"
		log.Println(sql)
		//log.Println(body)
		log.Println(id)
		stmt, err := config.GetReplicaDB(name).Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
		}
		_, err = stmt.Exec(string(body), idInt)
		//db.Close()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)

			log.Println(err.Error())
			return
		}
		stmt.Close()
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": "ok"})
		w.Write(response)
		return
	}
	//Db.Exec(initializeStr)
	//log.Print("Sqlite database: " + dbName)
	//sql := "SELECT \"DatasetName\",\"ItemId\",\"ItemInfo\",\"AdvancedDrawingInfo\" FROM \"GDB_ServiceItems\""
	sql := "SELECT " + fieldStr + " FROM " + config.Collector.Schema + config.DblQuote("GDB_ServiceItems") + " where " + config.DblQuote("OBJECTID") + "=?"
	log.Printf("Query: "+sql+"%v", idInt)

	stmt, err := config.GetReplicaDB(name).Prepare(sql)
	if err != nil {
		log.Println(err.Error())
		//w.Write([]byte(err.Error()))
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		w.Write(response)

		return
	}
	//rows := stmt.QueryRow(id)
	var itemInfo []byte
	err = stmt.QueryRow(idInt).Scan(&itemInfo)
	//rows, err := Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		w.Write(response)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(itemInfo)
	/*
		for rows.Next() {
			err = rows.Scan(&itemInfo)
			w.Header().Set("Content-Type", "application/json")

			w.Write(itemInfo)
			//fmt.Println(string(itemInfo))
		}
		rows.Close() //good habit to close
	*/
	//db.Close()

}
