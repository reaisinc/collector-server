package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	config "github.com/traderboy/collector-server/config"
	structs "github.com/traderboy/collector-server/structs"
)

func Updates(name string, id string, parentTableName string, tableName string, updateTxt string, globalIdName string, joinField string, parentObjectID string) []byte {
	//log.Println(updateTxt)
	//var updates structs.Record
	var updates []structs.Feature
	decoder := json.NewDecoder(strings.NewReader(updateTxt)) //r.Body

	err := decoder.Decode(&updates)
	if err != nil {
		panic(err)
	}
	//defer r.Body.Close()
	cols := ""
	sep := ""
	c := 1
	//var vals := []interface{}
	//objectid := 1
	var objectid int
	//var globalID string
	var results []interface{}
	//var objId int
	//don't update these fields
	//globaloidname,joinField,oidname

	for num, i := range updates {
		var vals []interface{}

		result := map[string]interface{}{}
		for key, j := range i.Attributes {
			//fmt.Println(key + ":  ")
			//var objectid = updates[0].Attributes["OBJECTID"]
			//var globalId = updates[0].Attributes["GlobalID"]
			if key == joinField { //"GlobalGUID" {
				continue
			}
			//never update GlobalID
			if key == "GlobalID" {
				continue
			}
			if key == parentObjectID {
				objectid = int(j.(float64))
				result["objectId"] = objectid

				//objId = c
				//c++
				//} else if key == "GlobalID" {
				//	globalID = j.(string)
				//	result["globalId"] = globalID
			} else {
				//if j != nil {
				//need to handle nulls
				if j == nil {
					cols += sep + config.DblQuote(key) + "=null"
				} else {
					cols += sep + config.DblQuote(key) + "=" + config.GetParam(config.Collector.DefaultDataSource, c)
					vals = append(vals, j)
					c++
				}
				sep = ","
				//fmt.Println(j)
				//}
			}
		}
		//cast(strftime('%s','now') as int)

		if config.Collector.Projects[name].Layers[id].EditFieldsInfo != nil {
			//joinField = config.Collector.Projects[name].Layers[id]["joinField"].(string)
			//for key, j := range config.Collector.Projects[name].Layers[id]["editFieldsInfo"] {
			current_time := time.Now().Local()
			//if rec, ok := config.Collector.Projects[name].Layers[id].EditFieldsInfo.(map[string]interface{}); ok {
			//if rec, ok := config.Collector.Projects[name].Layers[id].EditFieldsInfo.(map[string]interface{}); ok {
			//cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreatorField) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
			//vals = append(vals, config.Collector.Username)
			//p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
			//i.Attributes["creatorField"] = config.Collector.Username
			//c++
			cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].EditFieldsInfo.EditorField) + "=" + config.GetParam(config.Collector.DefaultDataSource, c) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
			vals = append(vals, config.Collector.Username)
			//p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
			i.Attributes["editorField"] = config.Collector.Username
			updates[num].Attributes["editorField"] = config.Collector.Username
			c++

			//p += sep + config.DbTimeStamp                                  //julianday('now')"
			//i.Attributes["creationDateField"] = current_time.Unix() * 1000 //DateToNumber(current_time.Year(), current_time.Month(), current_time.Day())
			//p += sep + config.DbTimeStamp                                  //julianday('now')"
			cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].EditFieldsInfo.EditDateField) + "=" + config.GetParam(config.Collector.DefaultDataSource, c) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
			vals = append(vals, current_time.Unix()*1000)
			i.Attributes["editDateField"] = current_time.Unix() * 1000 //DateToNumber(current_time.Year(), current_time.Month(), current_time.Day())
			updates[num].Attributes["editDateField"] = i.Attributes["editDateField"]
			c++

			/*
				for key, j := range rec {
					if key == "creatorField" || key == "editorField" {
						if key == "creatorField" {
							continue
						}
						vals = append(vals, config.Collector.Username)
						cols += sep + config.DblQuote(j.(string)) + "=" + config.GetParam(config.Collector.DefaultDataSource, c) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
						i.Attributes[key] = config.Collector.Username
						updates[num].Attributes[key] = config.Collector.Username
						c++
					} else if key == "creationDateField" || key == "editDateField" {
						//vals = append(vals, "julianday('now')")
						if key == "creationDateField" {
							continue
						}
						cols += sep + config.DblQuote(j.(string)) + "=" + config.DbTimeStamp // "=((julianday('now') - 2440587.5)*86400.0*1000)"
						//julianday('now')"
						i.Attributes[key] = current_time.Unix() * 1000
						//DateToNumber(current_time.Year(), current_time.Month(), current_time.Day())
						updates[num].Attributes[key] = i.Attributes[key]
					}
				}
			*/
			//}
			/*
				cols += sep + config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creatorField"]
				p += sep + config.GetParam(c)
				c++
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creatorField"] = config.Collector.Username
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["editorField"]=config.Collector.Username
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creationDateField"]=
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["editDateField"]
			*/
		}
		//add objectid last
		vals = append(vals, objectid)
		//tableName = strings.Replace(tableName, "_evw", "", -1)

		log.Println("update " + config.Collector.Schema + tableName + " set " + cols + " where " + config.DblQuote(parentObjectID) + "=" + config.GetParam(config.Collector.DefaultDataSource, len(vals)))
		log.Print(vals)
		//log.Print(objId)
		var sql string
		//if config.Collector.DefaultDataSource == structs.PGSQL {
		//	sql = "update " + config.Collector.Schema + config.DblQuote(tableName) + " set " + cols + " where " + config.DblQuote(parentObjectID) + "=" + config.GetParam(config.Collector.DefaultDataSource, len(vals))
		//} else if config.Collector.DefaultDataSource == structs.SQLITE3 {
		sql = "update " + tableName + " set " + cols + " where " + config.DblQuote(parentObjectID) + "=?"
		//}

		stmt, err := config.GetReplicaDB(name).Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		//err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
		_, err = stmt.Exec(vals...)
		if err != nil {
			log.Println(err.Error())
		}
		stmt.Close()
		result["success"] = true
		result["globalId"] = nil
		results = append(results, result)

		if i.Geometry != nil {
			//log.Println("Checking geometry")
			//var geometry string
			//geometry = getESRIPoint(i.Geometry.X, i.Geometry.Y, config.Collector.Projects[name].ReplicaPath)
			geometry := fmt.Sprintf("st_point('point(%v %v)',3857)", i.Geometry.X, i.Geometry.Y)
			//cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].ShapeFieldName) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
			//p += sep + geometry
			//vals = append(vals, geometry)
			//p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
			//i.Attributes["creatorField"] = config.Collector.Username
			objectidstr := strconv.Itoa(objectid)
			sql = "update " + config.Collector.Schema + tableName + " set " + config.Collector.Projects[name].Layers[id].ShapeFieldName + "=" + geometry + " where " + parentObjectID + "=" + objectidstr
			log.Println(sql)
			//convert ? to actual delimited values
			//need to update spatial index
			err := runSqliteCmd(sql, config.Collector.Projects[name].ReplicaPath)
			log.Println(err)
			//c++
		}

		/*
			select pos-1  from services,jsonb_array_elements(json->'features') with ordinality arr(elem,pos) where type='query' and layerId=0 and elem->'attributes'->>'OBJECTID'='$1')::int

			update services set json=jsonb_set(json,
			'{features,26,attributes}',
			'{"OBJECTID":27,"acres":3.12,"lease_site":0,"feature_type":1,"climatic_zone":2,"quad_name":"077-SE-196","elevation":6048,"permittee":"Lorraine / Elsie Begay","homesite_id":"H61A"}'::jsonb,
			false) where type='query' and layerId=0;
		*/
		//sql = "update services set json=jsonb_set(json, array('features',elem_index::text, ,false) from (select pos - 1 as elem_index from services,jsonb_array_elements(json->'features') with ordinality arr(elem,pos) where type='query' and layerId=0 and elem->'attributes'->>'OBJECTID'='$2')"

		updateTxt = updateTxt[15 : len(updateTxt)-2]
		if config.Collector.DefaultDataSource == structs.PGSQL {
			//update the same data in Postgresql
			sql = "update " + config.Collector.Schema + config.DblQuote(tableName) + " set " + cols + " where " + config.DblQuote(parentObjectID) + "=" + config.GetParam(config.Collector.DefaultDataSource, len(vals))
			stmt, err := config.Collector.DatabaseDB.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			//err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
			_, err = stmt.Exec(vals...)
			if err != nil {
				log.Println(err.Error())
			}
			stmt.Close()

			//now update the JSON
			sql = "select pos-1  from " + config.Collector.Schema + "services,jsonb_array_elements(json->'features') with ordinality arr(elem,pos) where type='query' and layerId=$1 and elem->'attributes'->>'OBJECTID'=$2"

			log.Println(sql)
			log.Printf("Layer ID: %v, ObjectID: %v\n", id, objectid)
			//log.Println(id)
			//log.Print("Objectid: ")
			//log.Println(objectid)
			rows, err := config.Collector.DatabaseDB.Query(sql, id, objectid)

			var rowId int
			for rows.Next() {
				err := rows.Scan(&rowId)
				if err != nil {
					log.Fatal(err)
				}
			}
			rows.Close()
			sql = "update " + config.Collector.Schema + "services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}',$1::jsonb,false) where type='query' and layerId=$2"
			log.Println(sql)
			stmt, err = config.Collector.DatabaseDB.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			log.Println(updateTxt)
			log.Println(id)
			_, err = stmt.Exec(updateTxt, id)
			if err != nil {
				log.Println(err.Error())
			}
			stmt.Close()

		} else if config.Collector.DefaultDataSource == structs.SQLITE3 {
			sql = "select json from services where type='query' and layerId=?"
			stmt, err = config.Collector.Configuration.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			rows, err := config.Collector.Configuration.Query(sql, id, objectid)

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
				//if i.Attributes[parentObjectID] != nil {
				oid := int(i.Attributes[parentObjectID].(float64))
				if oid == objectid {
					//i.Attributes["OBJECTID"]
					fieldObj.Features[k].Attributes = updates[num].Attributes
					break
				}
				//}
			}
			var jsonstr []byte
			jsonstr, err = json.Marshal(fieldObj)
			if err != nil {
				log.Println(err)
			}
			//log.Println(string(jsonstr))
			tx, err := config.Collector.Configuration.Begin()
			if err != nil {
				log.Fatal(err)
			}

			sql = "update " + config.Collector.Schema + "services set json=? where type='query' and layerId=?"
			log.Println(sql)

			stmt, err = tx.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}

			idInt, _ := strconv.Atoi(id)
			//log.Printf("%v\n%v", string(jsonstr), idInt)
			//sql = "PRAGMA synchronous = OFF;PRAGMA cache_size=100000;PRAGMA journal_mode=WAL;"
			//tx.Exec(sql)

			_, err = tx.Stmt(stmt).Exec(string(jsonstr), idInt)
			if err != nil {
				log.Println(err.Error())
			}
			tx.Commit()
			stmt.Close()
			//sql = "update services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}',$1::jsonb,false) where type='query' and layerId=$2"
		}
	}
	response, _ := json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": results, "deleteResults": []string{}})
	return response

	//curl -H "Content-Type: application/x-www-form-urlencoded" -X POST -d 'rollbackOnFailure=true&updates=[{"attributes":{"OBJECTID":3,"permittee":"Jack/Bessie Hatathlie","homesite_id":"9w3hdseq78dy","range_unit":551,"acres":3,"lease_site":0,"feature_type":0,"climatic_zone":2,"quad_name":"099-NW-004","elevation":6040,"permittee_globalid":"{D1A2F0B1-6F46-477A-80A9-CF550915B6BB}","has_permittee":1}}]&f=json' http://localhost:81/arcgis/rest/services/leasecompliance2016/FeatureServer/0/applyEdits

	//curl -H "Content-Type: application/x-www-form-urlencoded" -X POST -d 'rollbackOnFailure=true&adds=[{"geometry":"attributes":{"OBJECTID":3,"permittee":"Jack/Bessie Hatathlie","homesite_id":"9w3hdseq78dy","range_unit":551,"acres":3,"lease_site":0,"feature_type":0,"climatic_zone":2,"quad_name":"099-NW-004","elevation":6040,"permittee_globalid":"{D1A2F0B1-6F46-477A-80A9-CF550915B6BB}","has_permittee":1}}]&f=json' http://localhost:81/arcgis/rest/services/leasecompliance2016/FeatureServer/0/applyEdits

	//var jsonvals []interface{}
	//updateTxt := "[{\"attributes\":{\"OBJECTID\":27,\"acres\":3.15,\"lease_site\":0,\"feature_type\":1,\"climatic_zone\":2,\"quad_name\":\"077-SE-196\",\"elevation\":6048,\"permittee\":\"Lorraine / Elsie Begay\",\"homesite_id\":\"H61A\"}}]"
	//updateTxt = strings.Replace(updateTxt[15:len(updateTxt)-1], "\"", "\\\"", -1)

	//jsonvals = append(jsonvals, updateTxt)
	//jsonvals = append(jsonvals, id)
	//jsonvals = append(jsonvals, rowId)

	/*
		_, err = stmt.Exec(jsonvals...)
		if err != nil {
			log.Println(err.Error())
		}
	*/
	/*
		sql = "update services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}','" + updateTxt + "'::jsonb,false) where type='query' and layerId=$1"
		stmt, err = config.Collector.DatabaseDB.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = stmt.Exec(strconv.Atoi(id))
		if err != nil {
			log.Println(err.Error())
		}
	*/

	//log.Println(sql)
	//log.Println(jsonvals)
	/*
		_, err = stmt.Exec(sql, updateTxt, id)
		if err != nil {
			log.Println(err.Error())
		}
	*/

	/*
		var jsonvals []interface{}
		jsonvals = append(jsonvals, updateTxt)

		jsonvals = append(jsonvals, id)

	*/

	//find the matching OBJECTID in the query.json file and update fields and save back to disk
	/*
		for _, i := range updates {
			for _, j := range fields.Fields ["features"] {
				for _, k := range updates[i]["attributes"] {

				}

			}
		}
	*/

	/*
		err2 := json.Unmarshal(r.FormValue("updates"), &updates)
		if err2 != nil {
			log.Println("Error reading configuration file: " + r.FormValue("updates"))
			log.Println(err2.Error())
		}
	*/
	/*
	   decoder := json.NewDecoder(r.Body)
	       var t test_struct
	       err := decoder.Decode(&t)
	       if err != nil {
	           panic(err)
	       }
	       defer req.Body.Close()
	*/

	//var jsonFields=JSON.parse(file)
	//log.Println("sqlite: " + replicaDb)
	//var db = new sqlite3.Database(replicaDb)
	/*
		var sqlstr = "update " + outFields + " from " +
			config.Services[name]["relationships"][id]["dTable"].(string) +
			" where " +
			config.Services[name]["relationships"][id]["dJoinKey"].(string) + " in (select " +
			config.Services[name]["relationships"][id]["oJoinKey"].(string) + " from " +
			config.Services[name]["relationships"][id]["oTable"].(string) +
			" where OBJECTID=?)"

		db, err := sql.Open("sqlite3", replicaDb)
		if err != nil {
			log.Fatal(err)
		}
		defer Db.Close()
		stmt, err := Db.Prepare(sqlstr)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		//outArr := []interface{}{}
		rows, err := stmt.Query(id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		//final_result := map[int]map[string]string{}
		//works final_result := map[int]map[string]interface{}{}
		final_result := make([]interface{}, 0)
		result_id := 0
	*/

	//var updates = JSON.parse(req.body.updates)//JSON.parse(req.query.updates)
	/*
			var fs = require("fs')
			var path=DataPath+"/"+name") +"/FeatureServer."+id") + ".query.json"
		  var file = fs.readFileSync(path, "utf8")
		  var json=JSON.parse(file)
		  var results=[]
		  var fields=[]
		  var values=[]

		  for(var u=0;u<updates.length;u++)
		  {
			  for(var i=0;i<json.features.length;i++)
			  {
			  	//log.Println(json.features[i]['attributes']['OBJECTID'] + ":  " + updates[u].attributes['OBJECTID'])
			  	if(json.features[i]['attributes']['OBJECTID']==updates[u].attributes['OBJECTID'])
			  	{
			  		//json.features.[i]['attributes']=updates
			  		for(var j in updates[u].attributes)
			  		{
			  			for(var k in json.features[i]['attributes'])
			  			{
			  				if(j==k)
			  				{
			  					if(json.features[i]['attributes'][k] != updates[u].attributes[j])
			  					{
			  					    log.Println("Updating record: " + updates[u].attributes['OBJECTID'] + " " + k + "   values: " + json.features[i]['attributes'][k]+ " to " + updates[u].attributes[j] )
			  					    json.features[i]['attributes'][k]=updates[u].attributes[j]
		  	              fields.push(k+"=?")
		  	              values.push(updates[0].attributes[j])
			  					    break
			  				  }
			  				}
			  			}
			  		}
			  		results.push({"objectId":updates[u].attributes['OBJECTID'],"globalId":null,"success":true})
			  		break
			  	}
			  }
		  }
		  if(fields.length>0){
			  //search for id and update all fields
			  fs.writeFileSync(path, JSON.stringify(json), "utf8")

			  //now update the replica database

			  values.push(parseInt(id")))

			  var replicaDb = ReplicaPath + "/"+name")+".geodatabase"
			  log.Println("sqlite: " + replicaDb)
			  var db = new sqlite3.Database(replicaDb)
			  //create update statement from json
			  log.Println("UPDATE " + name") + " SET "+fields.join(",")+" WHERE OBJECTID = ?")
			  log.Println( values )

			  Db.run("UPDATE " + name") + " SET "+fields.join(",")+" WHERE OBJECTID = ?", values)
		  }else{
		 	  results={"objectId":updates.length>0?updates[0].attributes['OBJECTID']:0,"globalId":null,"success":true}
		 	}
	*/
	//update json file with updates
}
func UpdatesFile(name string, id string, parentTableName string, updatesTxt string, joinField string, globalIdName string, parentObjectID string) []byte {
	var updates structs.Record
	decoder := json.NewDecoder(strings.NewReader(updatesTxt)) //r.Body
	err := decoder.Decode(&updates)
	if err != nil {
		panic(err)
	}
	var objectid int
	//var globalID string
	var results []interface{}

	var fieldObj structs.FeatureTable
	current_time := time.Now().Local()

	jsonFile := config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json"
	//log.Println(jsonFile)
	file, err1 := ioutil.ReadFile(jsonFile)
	if err1 != nil {
		log.Println(err1)
	}
	err = json.Unmarshal(file, &fieldObj)
	if err != nil {
		log.Println("Error unmarshalling fields into features object: " + string(file))
		log.Println(err.Error())
	}

	//var objId int
	for k, i := range fieldObj.Features {
		//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
		//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))
		if int(i.Attributes[parentObjectID].(float64)) == objectid {
			//i.Attributes["OBJECTID"]
			fieldObj.Features[k].Attributes = updates[0].Attributes
			//if edit, save username and timestamp
			if config.Collector.Projects[name].Layers[id].EditFieldsInfo != nil {
				//fieldObj.Features[k].Attributes[config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreatorField] = config.Collector.Username
				fieldObj.Features[k].Attributes[config.Collector.Projects[name].Layers[id].EditFieldsInfo.EditorField] = config.Collector.Username
				//fieldObj.Features[k].Attributes[config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreationDateField] = current_time.Unix() * 1000
				fieldObj.Features[k].Attributes[config.Collector.Projects[name].Layers[id].EditFieldsInfo.EditDateField] = current_time.Unix() * 1000
			}

			break
		}
	}

	var jsonstr []byte
	jsonstr, err = json.Marshal(fieldObj)
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(jsonFile, jsonstr, 0644)
	if err != nil {
		log.Println(err1)
	}
	//write json back to file
	result := map[string]interface{}{}
	result["objectId"] = objectid
	result["success"] = true
	result["globalId"] = nil
	results = append(results, result)
	response, _ := json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": results, "deleteResults": []string{}})
	return response
}
