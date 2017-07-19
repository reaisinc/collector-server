package routes

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/traderboy/arcrestgo/structs"
	config "github.com/traderboy/collector-server/config"
)

func Deletes(name string, id string, parentTableName string, tableName string, deletesTxt string, globalIdName string, parentObjectID string) []byte {
	//deletesTxt should be a objectId
	var objectid, _ = strconv.Atoi(deletesTxt)
	var results []interface{}
	result := map[string]interface{}{}
	result["objectId"] = objectid
	result["success"] = true
	result["globalId"] = nil
	results = append(results, result)
	//delete from table
	log.Println("delete from " + config.Schema + tableName + " where " + config.DblQuote(parentObjectID) + " in (" + config.GetParam(1) + ")")
	log.Println("delete objectids:  " + deletesTxt + "/" + strconv.Itoa(objectid))
	var sql = "delete from " + config.Schema + tableName + " where " + config.DblQuote(parentObjectID) + " in (" + config.GetParam(1) + ")"
	stmt, err := config.DbQuery.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	//err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
	_, err = stmt.Exec(objectid)
	if err != nil {
		log.Println(err.Error())
	}
	stmt.Close()

	if config.DbSource == config.PGSQL {
		sql := "select pos-1  from " + config.Schema + "services,jsonb_array_elements(json->'features') with ordinality arr(elem,pos) where type='query' and layerId=$1 and elem->'attributes'->>'OBJECTID'=$2"

		log.Println(sql)
		log.Printf("Layer ID: %v", id)
		log.Printf("Objectid: %v", objectid)

		rows, err := config.Db.Query(sql, id, objectid)

		var rowId int
		for rows.Next() {
			err := rows.Scan(&rowId)
			if err != nil {
				log.Fatal(err)
			}
		}
		rows.Close()
		//sql = "update services set json=json->'features' - " + strconv.Itoa(rowId) + " where type='query' and layerId=$1"
		sql = "update " + config.Schema + "services set json=json #- '{features," + strconv.Itoa(rowId) + "}' where type='query' and layerId=$1"
		log.Println(sql)
		log.Printf("Row id: %v", rowId)
		stmt, err := config.Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			log.Println(err.Error())
		}
		stmt.Close()

	} else if config.DbSource == config.SQLITE3 {
		sql := "select json from services where type='query' and layerId=?"
		stmt, err := config.Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		rows, err := config.Db.Query(sql, id, objectid)

		var row []byte
		for rows.Next() {
			err := rows.Scan(&row)
			if err != nil {
				log.Fatal(err)
			}
		}
		rows.Close()
		stmt.Close()

		var fieldObj structs.FeatureTable
		//map[string]map[string]map[string]
		err = json.Unmarshal(row, &fieldObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(row))
			log.Println(err.Error())
		}
		for k, i := range fieldObj.Features {
			//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
			//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))
			if int(i.Attributes[parentObjectID].(float64)) == objectid {
				//i.Attributes["OBJECTID"]
				//fieldObj.Features = fieldObj.Features[k]
				fieldObj.Features = append(fieldObj.Features[:k], fieldObj.Features[k+1:]...)
				//fieldObj.Features[k].Attributes = updates[num].Attributes
				break
			}
		}
		var jsonstr []byte
		jsonstr, err = json.Marshal(fieldObj)
		if err != nil {
			log.Println(err)
		}
		tx, err := config.Db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		sql = "update " + config.Schema + "services set json=? where type='query' and layerId=?"

		stmt, err = tx.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		idInt, _ := strconv.Atoi(id)

		_, err = tx.Stmt(stmt).Exec(string(jsonstr), idInt)
		if err != nil {
			log.Println(err.Error())
		}
		tx.Commit()
		stmt.Close()
		//sql = "update services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}',$1::jsonb,false) where type='query' and layerId=$2"
	}
	response, _ := json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": []string{}, "deleteResults": results})
	return response

}
