package routes

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	config "github.com/traderboy/collector-server/config"
)
func xml(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		//w.Header().Set("Content-Type", "text/xml")
		body := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + r.URL.Query().Get("xml")
		w.Write([]byte(body))
	}
	
func xml_id(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	//idInt, _ := strconv.Atoi(id)
	dbPath := r.URL.Query().Get("db")
	tableName := config.Project.Services[name]["layers"][id]["data"].(string)
	tableName = strings.ToUpper(tableName)

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/xml/" + id)
	var dbName = config.ReplicaPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
	if len(dbPath) > 0 {
		if config.DbSource != config.PGSQL {
			if config.DbSqliteDbName != dbPath {
				if config.DbSqliteQuery != nil {
					config.DbSqliteQuery.Close()
				}
				config.DbSqliteQuery = nil
			}
			config.DbSqliteDbName = dbPath
			dbName = "file:" + dbPath + config.SqlWalFlags //+ "?PRAGMA journal_mode=WAL"
		}
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
	if config.DbSource == config.PGSQL {
		config.DbSqliteQuery = config.DbQuery
	} else {
		if config.DbSqliteQuery == nil {
			//config.DbSqliteQuery, err = sql.Open("sqlite3", "file:"+dbName+"?PRAGMA journal_mode=WAL")
			config.DbSqliteQuery, err = sql.Open("sqlite3", dbName)
			if err != nil {
				log.Fatal(err)
			}
		}
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
		sql := "update " + config.Schema + config.DblQuote("GDB_Items") + " set " + config.DblQuote("Definition") + "=? where " + config.DblQuote("PhysicalName") + "=?" //OBJECTID=?"
		stmt, err := config.DbSqliteQuery.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)

			return
		}
		_, err = stmt.Exec(body, tableName)
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
	//Db.Exec(initializeStr)
	log.Print("Sqlite database: " + dbName)
	//sql := "SELECT \"DatasetName\",\"ItemId\",\"ItemInfo\",\"AdvancedDrawingInfo\" FROM \"GDB_ServiceItems\""
	sql := "SELECT " + config.DblQuote("Definition") + " FROM " + config.Schema + config.DblQuote("GDB_Items") + " where " + config.DblQuote("PhysicalName") + "=?" //OBJECTID=?"
	log.Printf("Query: "+sql+"%v", tableName)

	stmt, err := config.DbSqliteQuery.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		w.Write(response)
	}
	//rows := stmt.QueryRow(id)
	var itemInfo []byte
	err = stmt.QueryRow(tableName).Scan(&itemInfo)
	//rows, err := Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		w.Write(response)

		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(itemInfo)
}
