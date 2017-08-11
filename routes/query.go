package routes

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	config "github.com/traderboy/collector-server/config"
	structs "github.com/traderboy/collector-server/structs"
)

func queryRelatedRecordsDB(name string, id string, relationshipId string, objectIds string, objectId int, outFields string, parentObjectID string, dID int) []byte {
	var sql string
	var fields []byte
	var fieldsArr []structs.Field
	//var outFields string

	if config.Collector.DefaultDataSource == structs.PGSQL {
		sql = "select json->'fields' from " + config.Collector.Schema + "services where service=$1 and name=$2 and layerid=$3 and type=$4"
		log.Printf("select json->'fields' from "+config.Collector.Schema+"services where service='%v' and name='%v' and layerid=%v and type='%v'", name, "FeatureServer", dID, "")
		stmt, err := config.Collector.DatabaseDB.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		err = stmt.QueryRow(name, "FeatureServer", dID, "").Scan(&fields)
		if err != nil {
			log.Println(err.Error())
			//w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return []byte("{\"fields\":[],\"relatedRecordGroups\":[]}")
		}
		err = json.Unmarshal(fields, &fieldsArr)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		/*
			var outFieldsArr []string
			if outFields != "*" {
				outFieldsArr = strings.Split(outFields, ",")
			}
		*/
		outFields = ""
		pre := ""
		//need to change date fields to TO_CHAR(created_date, 'J')
		for _, i := range fieldsArr {
			//log.Println("%v %v\n", k, i)
			if i.Type == "esriFieldTypeDate" {
				//outFields += pre + "TO_CHAR(" + i.Name + ", 'J') as " + i.Name
				outFields += pre + "(CAST (to_char(" + i.Name + ", 'J') AS INT) - 2440587.5)*86400.0*1000  as " + i.Name
			} else {
				outFields += pre + config.DblQuote(i.Name)
			}
			pre = ","
			//outFields += config.DblQuote(fieldObj.Features[k].Attributes)
			//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
			//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))
			//if int(i.Attributes[parentObjectID].(float64)) == objectid {
			//i.Attributes["OBJECTID"]
			//fieldObj.Features[k].Attributes = updates[0].Attributes
			//break
			//}
		}
		//log.Println("%v", outFieldsArr)

	} else if config.Collector.DefaultDataSource == structs.SQLITE3 {
		sql = "select json from services where service=? and name=? and layerid=? and type=?"
		log.Printf("select json from services where service='%v' and name='%v' and layerid=%v and type='%v'", name, "FeatureServer", dID, "")
		stmt, err := config.Collector.Configuration.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		err = stmt.QueryRow(name, "FeatureServer", dID, "").Scan(&fields)
		if err != nil {
			log.Println(err.Error())
			//w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return []byte("{\"fields\":[],\"relatedRecordGroups\":[]}")
		}
		//fields = fields["fields"]

		var fieldObj structs.FeatureTable
		//map[string]map[string]map[string]
		err = json.Unmarshal(fields, &fieldObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		fieldsArr = fieldObj.Fields

	}
	//Fields            []Field   `json:"fields,omitempty"`
	//create outFields

	//fields = fields["fields"]
	//map[string]map[string]map[string]

	//for

	//_, err = w.Write(fields)
	//return
	//var replicaDb = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
	//var tableName = config.Collector.Projects[name].Relationships[relationshipId]["dTable"].(string)

	//log.Println(tableName)
	//var layerId = int(config.Services[name].Relationships[relationshipId]["dId"].(float64))
	//var jsonFields=JSON.parse(file)
	//log.Println("sqlite: " + replicaDb)
	//var db = new sqlite3.Database(replicaDb)
	joinField := config.Collector.Projects[name].Relationships[relationshipId].OJoinKey
	//if joinField == "GlobalID" || joinField == "GlobalGUUD" {
	//	joinField = "substr(" + joinField + ", 2, length(" + joinField + ")-2)"
	//}
	var sqlstr = "select " + outFields + " from " + config.Collector.Schema +
		config.DblQuote(config.Collector.Projects[name].Relationships[relationshipId].DTable) +
		" where " +
		config.DblQuote(config.Collector.Projects[name].Relationships[relationshipId].DJoinKey) +
		" in (select " +
		config.DblQuote(joinField) + " from " +
		config.Collector.Schema + config.DblQuote(config.Collector.Projects[name].Relationships[relationshipId].OTable) +
		" where " + config.DblQuote(parentObjectID) + " in(" + config.GetParam(config.Collector.DefaultDataSource, 1) + "))"

	//_, err = w.Write([]byte(sqlstr))
	log.Println(strings.Replace(sqlstr, config.GetParam(config.Collector.DefaultDataSource, 1), objectIds, -1))

	stmt, err := config.GetReplicaDB(name).Prepare(sqlstr)
	if err != nil {
		log.Fatal(err)
	}

	//outArr := []interface{}{}
	//relationshipIdInt, _ := strconv.Atoi(relationshipId)
	objectidArr, _ := strconv.Atoi(objectIds)
	rows, err := stmt.Query(objectidArr) //relationshipIdInt
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//var colLookup = map[string]interface{}{"objectid": "OBJECTID", "globalid": "GlobalID", "creationdate": "CreationDate", "creator": "Creator", "editdate": "EditDate", "editor": "Editor"}
	var colLookup = map[string]string{"objectid": "OBJECTID", "globalguid": "GlobalGUID", "globalid": "GlobalID", "creationdate": "CreationDate", "creator": "Creator", "editdate": "EditDate", "editor": "Editor", "comments": "Comments"}
	var guuids = map[string]int{"GlobalGUID": 1, "GlobalID": 1}
	var dates = map[string]int{"created_date": 1, "last_edited_date": 1}
	columns, _ := rows.Columns()
	//colTypes, _ := rows.ColumnTypes()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//for i, col := range colTypes {
	//	log.Printf("%v: %v", col.Name, col.DatabaseTypeName)
	//}
	//final_result := map[int]map[string]string{}
	//works final_result := map[int]map[string]interface{}{}
	final_result := make([]interface{}, 0)
	result_id := 0
	//log.Println("Query ran successfully")
	for rows.Next() {
		//log.Println("next row")
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
			//log.Println(i)
		}
		rows.Scan(valuePtrs...)
		//tmp_struct := map[string]string{}
		tmp_struct := map[string]interface{}{}

		for i, col := range columns {
			//var v interface{}
			val := values[i]

			if colLookup[col] != "" {
				col = colLookup[col]
			}
			//fmt.Printf("Integer: %v=%v\n", col, val)
			switch t := val.(type) {
			case int:
				//fmt.Printf("Integer: %v=%v\n", col, t)
				tmp_struct[col] = val
			case float64:
				tmp_struct[col] = val
				if dates[col] == 1 && val != nil {
					tmp_struct[col] = int(val.(float64))
				} else {
					tmp_struct[col] = val
				}
				//fmt.Printf("Float64: %v %v\n", col, val)
			case []uint8:
				//fmt.Printf("Col: %v (uint8): %v\n", col, t)
				b, _ := val.([]byte)
				tmp_struct[col] = fmt.Sprintf("%s", b)
				//sqlite
				if guuids[col] == 1 && tmp_struct[col] != nil {
					tmp_struct[col] = strings.Trim(tmp_struct[col].(string), "{}")
				}
				//fmt.Printf("Col: %v (uint8): %v\n", col, tmp_struct[col])

			case int64:
				//fmt.Printf("Integer 64: %v\n", t)
				tmp_struct[col] = val
			case string:
				tmp_struct[col] = fmt.Sprintf("%s", val)
				//pg
				if guuids[col] == 1 && tmp_struct[col] != nil {
					tmp_struct[col] = strings.Trim(tmp_struct[col].(string), "{}")
				}
				//fmt.Printf("String: %v=%v:  %v\n", col, val, tmp_struct[col])
			case bool:
				//fmt.Printf("Bool: %v\n", t)
				tmp_struct[col] = val
			case []interface{}:
				for i, n := range t {
					fmt.Printf("Item: %v= %v\n", i, n)
				}
			default:
				var r = reflect.TypeOf(t)
				tmp_struct[col] = r
				//fmt.Printf("Other:%v=%v\n", col, r)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		//log.Println(tmp_struct)
		record := map[string]interface{}{"attributes": tmp_struct}
		final_result = append(final_result, record)
		result_id++
	}
	//log.Println("Query end successfully")
	var resp []byte
	if len(final_result) > 0 {
		var result = map[string]interface{}{}
		//result["objectId"] = objectIds //strconv.Atoi(objectIds)
		//OBS! must convert objectID to int or it fails on Android
		oid, _ := strconv.Atoi(objectIds)
		result["objectId"] = oid
		result["relatedRecords"] = final_result
		resp, _ = json.Marshal(map[string]interface{}{"relatedRecordGroups": []interface{}{result}})
		resp = resp[1:]
	} else {
		resp = []byte("\"relatedRecordGroups\":[]}")
	}
	//convert fields to string
	fields, err = json.Marshal(fieldsArr)
	if err != nil {
		log.Println(err)
	}
	/*
		w.Write([]byte("{\"fields\":"))
		w.Write(fields)
		w.Write([]byte(","))
	*/
	response := append([]byte("{\"fields\":"), fields...)
	response = append(response, []byte(",")...)
	response = append(response, resp...)
	return response

}
func queryRelatedRecordsFile(name string, id string, relationshipId string, objectIds string, objectId int, outFields string, parentObjectID string, dID int) []byte {
	//have to find the joinAttribute value for source and destination
	/*
		var sqlstr = "select " + outFields + " from " + config.Collector.Schema +
			config.Collector.Projects[name].Relationships[relationshipId]["dTable"].(string) +
			" where " +
			config.Collector.Projects[name].Relationships[relationshipId]["dJoinKey"].(string) + " in (select " +
			config.Collector.Projects[name].Relationships[relationshipId]["oJoinKey"].(string) + " from " +
			config.Collector.Projects[name].Relationships[relationshipId]["oTable"].(string) +
			" where OBJECTID in(" + config.GetParam(1) + "))"
	*/
	var dJoinKey = config.Collector.Projects[name].Relationships[relationshipId].DJoinKey
	var oJoinKey = config.Collector.Projects[name].Relationships[relationshipId].OJoinKey

	jsonFile := fmt.Sprint(config.Collector.DataPath, string(os.PathSeparator), name+string(os.PathSeparator), "services", string(os.PathSeparator), "FeatureServer.", id, ".query.json")
	log.Println(jsonFile)
	file, err1 := ioutil.ReadFile(jsonFile)
	if err1 != nil {
		log.Println(err1)
	}
	var srcObj structs.FeatureTable

	//map[string]map[string]map[string]
	err := json.Unmarshal(file, &srcObj)
	if err != nil {
		log.Println("Error unmarshalling fields into features object: " + string(file))
		log.Println(err.Error())
	}

	var oJoinVal interface{}
	for _, i := range srcObj.Features {
		//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
		//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))

		if int(i.Attributes[parentObjectID].(float64)) == objectId {
			oJoinVal = i.Attributes[oJoinKey]
			//i.Attributes["OBJECTID"]
			//fieldObj.Features[k].Attributes = updates[num].Attributes
			break
			//record.RelatedRecord = append(record.RelatedRecord, fieldObj.Features[k].Attributes)
		}
	}
	//oJoinVal = strings.Replace(oJoinVal.(string), "{", "", -1)
	//oJoinVal = strings.Replace(oJoinVal.(string), "}", "", -1)
	//oJoinVal = strings.ToLower(oJoinVal.(string))

	//strconv.Itoa(int(dID.(float64)))
	jsonFile = fmt.Sprint(config.Collector.DataPath, string(os.PathSeparator), name, string(os.PathSeparator), "services", string(os.PathSeparator), "FeatureServer.", dID, ".query.json")
	log.Println(jsonFile)
	file, err1 = ioutil.ReadFile(jsonFile)
	if err1 != nil {
		log.Println(err1)
	}
	var fieldObj structs.FeatureTable

	//map[string]map[string]map[string]
	err = json.Unmarshal(file, &fieldObj)
	if err != nil {
		log.Println("Error unmarshalling fields into features object: " + string(file))
		log.Println(err.Error())
	}
	var relRecords structs.RelatedRecords
	relRecords.Fields = fieldObj.Fields

	var recordGroup structs.RelatedRecordGroup
	recordGroup.ObjectId = objectId

	//records.RelatedRecordGroups.ObjectId = objectId
	//records.ObjectId = objectId
	//records.RelatedRecord = map[string]interface{}
	//c := 0
	//log.Printf("Finding: %v", oJoinVal)

	for k, i := range fieldObj.Features {
		//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
		//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))

		if i.Attributes[dJoinKey] == oJoinVal {
			//if strings.EqualFold(i.Attributes[dJoinKey],oJoinVal)
			//log.Printf("Found: %v", i.Attributes[dJoinKey])
			var rec structs.RelatedRecord
			//i.Attributes["OBJECTID"]
			//fieldObj.Features[k].Attributes = updates[num].Attributes
			//break
			//var attributes structs.Attribute
			//attributes = fieldObj.Features[k].Attributes
			//rec.Attributes = append(rec.Attributes, fieldObj.Features[k].Attributes)
			rec.Attributes = fieldObj.Features[k].Attributes
			recordGroup.RelatedRecords = append(recordGroup.RelatedRecords, rec)
			//c++
		}

	}

	var jsonstr []byte
	//if c == 0 {
	//	records.RelatedRecordGroups = records.RelatedRecordGroups[:0]
	//}
	if len(recordGroup.RelatedRecords) > 0 {
		relRecords.RelatedRecordGroups = append(relRecords.RelatedRecordGroups, recordGroup)
	} else {
		relRecords.RelatedRecordGroups = make([]structs.RelatedRecordGroup, 0)
	}
	jsonstr, err = json.Marshal(relRecords)
	if err != nil {
		log.Println(err)
	}
	return jsonstr

	/*
		tx, err := config.Collector.DatabaseDB.Begin()
		if err != nil {
			log.Fatal(err)
		}

		var response []byte
		if len(final_result) > 0 {
			var result = map[string]interface{}{}
			result["objectId"] = objectIds //strconv.Atoi(objectIds)
			result["relatedRecords"] = final_result
			response, _ = json.Marshal(map[string]interface{}{"relatedRecordGroups": []interface{}{result}})
			response = response[1:]
		} else {
			response = []byte("\"relatedRecordGroups\":[]}")
		}
	*/

}

func queryDB(name string, id string, where string, outFields string, returnIdsOnly string, objectIds string) []byte {
	//if(req.query.outFields=='OBJECTID'){

	//idInt, _ := strconv.Atoi(id)
	//dbPath := r.URL.Query().Get("db")

	var objectIDName = config.Collector.Projects[name].Layers[id].Oidname
	var tableName = config.Collector.Projects[name].Layers[id].Data

	//returnGeometry := r.FormValue("returnGeometry")

	//log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query")
	//returnIdsOnly = true
	var sql string
	var err error
	var fields []byte
	var fieldsArr []structs.Field
	//var db *sql.DB
	//var rows *sql.Rows

	if config.Collector.DefaultDataSource == structs.PGSQL {
		sql = "select json->'fields' from " + config.Collector.Schema + "services where service=$1 and name=$2 and layerid=$3 and type=$4"
		log.Printf("select json->'fields' from "+config.Collector.Schema+"services where service='%v' and name='%v' and layerid=%v and type='%v'", name, "FeatureServer", id, "")
		stmt, err := config.Collector.DatabaseDB.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		err = stmt.QueryRow(name, "FeatureServer", id, "").Scan(&fields)
		if err != nil {
			log.Println(err.Error())
			//w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return []byte("{\"fields\":[],\"relatedRecordGroups\":[]}")
		}
		err = json.Unmarshal(fields, &fieldsArr)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		/*
			var outFieldsArr []string
			if outFields != "*" {
				outFieldsArr = strings.Split(outFields, ",")
			}
		*/
		outFields = ""
		pre := ""
		//need to change date fields to TO_CHAR(created_date, 'J')
		for _, i := range fieldsArr {
			//log.Println("%v %v\n", k, i)
			if i.Type == "esriFieldTypeDate" {
				//outFields += pre + "TO_CHAR(" + i.Name + ", 'J') as " + i.Name
				outFields += pre + "(CAST (to_char(" + i.Name + ", 'J') AS INT) - 2440587.5)*86400.0*1000  as " + i.Name
			} else {
				outFields += pre + config.DblQuote(i.Name)
			}
			pre = ","
			//outFields += config.DblQuote(fieldObj.Features[k].Attributes)
			//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
			//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))
			//if int(i.Attributes[parentObjectID].(float64)) == objectid {
			//i.Attributes["OBJECTID"]
			//fieldObj.Features[k].Attributes = updates[0].Attributes
			//break
			//}
		}
		//log.Println("%v", outFieldsArr)

	} else if config.Collector.DefaultDataSource == structs.SQLITE3 {
		sql = "select json from services where service=? and name=? and layerid=? and type=?"
		log.Printf("select json from services where service='%v' and name='%v' and layerid=%v and type='%v'", name, "FeatureServer", id, "")
		stmt, err := config.Collector.Configuration.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		err = stmt.QueryRow(name, "FeatureServer", id, "").Scan(&fields)
		if err != nil {
			log.Println(err.Error())
			//w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return []byte("{\"fields\":[],\"relatedRecordGroups\":[]}")
		}
		//fields = fields["fields"]

		var fieldObj structs.FeatureTable
		//map[string]map[string]map[string]
		err = json.Unmarshal(fields, &fieldObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		fieldsArr = fieldObj.Fields

	}

	//log.Println(r.FormValue("returnGeometry"))
	//log.Println(r.FormValue("outFields"))
	//sql := "select "+outFields + " from " +
	//where = ""

	//construct sql string
	if len(where) > 0 {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/where=" + where)
		sql = "select " + outFields + " from " + tableName + " where " + where
		//response := config.GetArcQuery(name, "FeatureServer", idInt, "query",objectIds,where)
		//w.Header().Set("Content-Type", "application/json")
		//var response = []byte("{\"objectIdFieldName\":\"OBJECTID\",\"globalIdFieldName\":\"GlobalID\",\"geometryProperties\":{\"shapeAreaFieldName\":\"Shape__Area\",\"shapeLengthFieldName\":\"Shape__Length\",\"units\":\"esriMeters\"},\"features\":[]}")
		//var response = []byte(`{"objectIdFieldName":"OBJECTID","globalIdFieldName":"GlobalID","geometryProperties":{"shapeLengthFieldName":"","units":"esriMeters"},"features":[]}`)
		//w.Write(response)

	} else if returnIdsOnly == "true" {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/objectids")
		sql = "select " + objectIDName + " from " + tableName //+ " where " + where

		/*
			response := config.GetArcService(name, "FeatureServer", idInt, "objectids", dbPath)
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".objectids.json")
				http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".objectids.json")
			}
		*/
	} else if len(objectIds) > 0 {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/objectIds=" + objectIds)
		sql = "select " + outFields + " from " + tableName + " where " + config.DblQuote(objectIDName) + " in (" + objectIds + ")"

		//only get the select objectIds
		//response := config.GetArcService(name, "FeatureServer", idInt, "query")
		/*
			response := config.GetArcQuery(name, "FeatureServer", idInt, "query", parentObjectID, objectIds, dbPath)

			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
				http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")

			}
		*/
		//if returnGeometry == "false" &&
	} else if strings.Index(outFields, objectIDName) > -1 { //r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/outfields=" + outFields)
		sql = "select " + outFields + " from " + tableName + " where " + config.DblQuote(objectIDName) + " in (" + objectIds + ")"

		/*
			response := config.GetArcService(name, "FeatureServer", idInt, "outfields", dbPath)
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".outfields.json")
				http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".outfields.json")
			}
		*/
	} else {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/else")
		sql = "select " + outFields + " from " + tableName
		/*
			response := config.GetArcService(name, "FeatureServer", idInt, "query", dbPath)
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
				http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")
			}
		*/
	}

	//http.ServeFile(w, r, config.Collector.DataPath + "/" + id  + "query.json")

	//var outFields string

	//Fields            []Field   `json:"fields,omitempty"`
	//create outFields

	//fields = fields["fields"]
	//map[string]map[string]map[string]

	//for

	//_, err = w.Write(fields)
	//return
	//var replicaDb = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
	//var tableName = config.Collector.Projects[name].Relationships[relationshipId]["dTable"].(string)

	//log.Println(tableName)
	//var layerId = int(config.Services[name].Relationships[relationshipId]["dId"].(float64))
	//var jsonFields=JSON.parse(file)
	//log.Println("sqlite: " + replicaDb)
	//var db = new sqlite3.Database(replicaDb)
	/*
		joinField := config.Collector.Projects[name].Layers[id].Oidname
		//if joinField == "GlobalID" || joinField == "GlobalGUUD" {
		//	joinField = "substr(" + joinField + ", 2, length(" + joinField + ")-2)"
		//}
		var sqlstr = "select " + outFields + " from " + config.Collector.Schema +
			config.DblQuote(config.Collector.Projects[name].Layers[id].DTable) +
			" where " +
			config.DblQuote(config.Collector.Projects[name].Layers[id].DJoinKey) +
			" in (select " +
			config.DblQuote(joinField) + " from " +
			config.Collector.Schema + config.DblQuote(config.Collector.Projects[name].Layers[id].OTable) +
			" where " + config.DblQuote(parentObjectID) + " in(" + config.GetParam(config.Collector.DefaultDataSource, 1) + "))"
	*/
	//_, err = w.Write([]byte(sqlstr))
	//log.Println(strings.Replace(sqlstr, config.GetParam(config.Collector.DefaultDataSource, 1), objectIds, -1))

	stmt, err := config.GetReplicaDB(name).Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	//outArr := []interface{}{}
	//relationshipIdInt, _ := strconv.Atoi(relationshipId)

	//if len(objectIds) > 0 {
	//objectidArr, _ := strconv.Atoi(objectIds)
	rows, err1 := stmt.Query() //relationshipIdInt
	//} else {
	//	rows, err := stmt.Query(nil)
	//}

	if err1 != nil {
		log.Fatal(err1)
	}
	defer rows.Close()

	//var colLookup = map[string]interface{}{"objectid": "OBJECTID", "globalid": "GlobalID", "creationdate": "CreationDate", "creator": "Creator", "editdate": "EditDate", "editor": "Editor"}
	var colLookup = map[string]string{"objectid": "OBJECTID", "globalguid": "GlobalGUID", "globalid": "GlobalID", "creationdate": "CreationDate", "creator": "Creator", "editdate": "EditDate", "editor": "Editor", "comments": "Comments"}
	var guuids = map[string]int{"GlobalGUID": 1, "GlobalID": 1}
	var dates = map[string]int{"created_date": 1, "last_edited_date": 1}
	columns, _ := rows.Columns()
	//colTypes, _ := rows.ColumnTypes()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//for i, col := range colTypes {
	//	log.Printf("%v: %v", col.Name, col.DatabaseTypeName)
	//}
	//final_result := map[int]map[string]string{}
	//works final_result := map[int]map[string]interface{}{}
	final_result := make([]interface{}, 0)
	result_id := 0
	//log.Println("Query ran successfully")
	for rows.Next() {
		//log.Println("next row")
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
			//log.Println(i)
		}
		rows.Scan(valuePtrs...)
		//tmp_struct := map[string]string{}
		tmp_struct := map[string]interface{}{}

		for i, col := range columns {
			//var v interface{}
			val := values[i]

			if colLookup[col] != "" {
				col = colLookup[col]
			}
			//fmt.Printf("Integer: %v=%v\n", col, val)
			switch t := val.(type) {
			case int:
				//fmt.Printf("Integer: %v=%v\n", col, t)
				tmp_struct[col] = val
			case float64:
				tmp_struct[col] = val
				if dates[col] == 1 && val != nil {
					tmp_struct[col] = int(val.(float64))
				} else {
					tmp_struct[col] = val
				}
				//fmt.Printf("Float64: %v %v\n", col, val)
			case []uint8:
				//fmt.Printf("Col: %v (uint8): %v\n", col, t)
				b, _ := val.([]byte)
				tmp_struct[col] = fmt.Sprintf("%s", b)
				//sqlite
				if guuids[col] == 1 && tmp_struct[col] != nil {
					tmp_struct[col] = strings.Trim(tmp_struct[col].(string), "{}")
				}
				//fmt.Printf("Col: %v (uint8): %v\n", col, tmp_struct[col])

			case int64:
				//fmt.Printf("Integer 64: %v\n", t)
				tmp_struct[col] = val
			case string:
				tmp_struct[col] = fmt.Sprintf("%s", val)
				//pg
				if guuids[col] == 1 && tmp_struct[col] != nil {
					tmp_struct[col] = strings.Trim(tmp_struct[col].(string), "{}")
				}
				//fmt.Printf("String: %v=%v:  %v\n", col, val, tmp_struct[col])
			case bool:
				//fmt.Printf("Bool: %v\n", t)
				tmp_struct[col] = val
			case []interface{}:
				for i, n := range t {
					fmt.Printf("Item: %v= %v\n", i, n)
				}
			default:
				var r = reflect.TypeOf(t)
				tmp_struct[col] = r
				//fmt.Printf("Other:%v=%v\n", col, r)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		//log.Println(tmp_struct)
		record := map[string]interface{}{"attributes": tmp_struct}
		final_result = append(final_result, record)
		result_id++
	}
	//log.Println("Query end successfully")
	var resp []byte
	if len(final_result) > 0 {
		var result = map[string]interface{}{}
		//result["objectId"] = objectIds //strconv.Atoi(objectIds)
		//OBS! must convert objectID to int or it fails on Android
		oid, _ := strconv.Atoi(objectIds)
		result["objectId"] = oid
		result["features"] = final_result
		resp, _ = json.Marshal(map[string]interface{}{"features": []interface{}{result}})
		resp = resp[1:]
	} else {
		resp = []byte("\"features\":[]}")
	}
	//convert fields to string
	fields, err = json.Marshal(fieldsArr)
	if err != nil {
		log.Println(err)
	}
	/*
		w.Write([]byte("{\"fields\":"))
		w.Write(fields)
		w.Write([]byte(","))
	*/
	response := append([]byte("{\"fields\":"), fields...)
	response = append(response, []byte(",")...)
	response = append(response, resp...)
	return response
}
