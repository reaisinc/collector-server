package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	config "github.com/traderboy/collector-server/config"
	structs "github.com/traderboy/collector-server/structs"
	"github.com/twinj/uuid"
)

func Adds(name string, id string, parentTableName string, tableName string, addsTxt string, joinField string, globalIdName string, parentObjectID string) []byte {
	var results []interface{}
	var objectid int
	var uuidstr string

	//log.Println(addsTxt)
	var adds []structs.Feature
	decoder := json.NewDecoder(strings.NewReader(addsTxt)) //r.Body
	err := decoder.Decode(&adds)
	if err != nil {
		panic(err)
	}
	cols := ""
	p := ""

	c := 1

	//need to update "GDB_RowidGenerators"->"base_id" to new id
	sql := "select \"base_id\"," + config.Collector.UUID + " from " + config.Collector.Schema + "\"GDB_RowidGenerators\" where \"registration_id\" in ( SELECT \"registration_id\" FROM " + config.Collector.Schema + "\"GDB_TableRegistry\" where \"table_name\"='" + parentTableName + "')"
	//sql := "select max(" + parentObjectID + ")+1," + config.Collector.UUID + " from " + tableName
	log.Println(sql)
	rows, err := config.GetReplicaDB(name).Query(sql)
	//defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&objectid, &uuidstr)
		if err != nil {
			//log.Fatal(err)
			objectid = 1
			uuidstr = strings.ToUpper(uuid.Formatter(uuid.NewV4(), uuid.FormatCanonicalCurly))
		}
	}
	defer rows.Close()
	sql = "update " + config.Collector.Schema + "\"GDB_RowidGenerators\" set \"base_id\"=" + (strconv.Itoa(objectid + 1)) + " where \"registration_id\" in ( SELECT \"registration_id\" FROM " + config.Collector.Schema + "\"GDB_TableRegistry\" where \"table_name\"='" + parentTableName + "')"
	log.Println(sql)
	_, err = config.GetReplicaDB(name).Exec(sql)
	if err != nil {
		log.Println(err.Error())
		//w.Write([]byte(err.Error()))
		//w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		//w.Write(response)
		return response
	}

	//var globalId string
	for _, i := range adds {
		var vals []interface{}
		sep := ""

		for key, j := range i.Attributes {
			if key == parentObjectID {
				i.Attributes[parentObjectID] = objectid
				cols += sep + config.DblQuote(key)
				p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
				sep = ","
				vals = append(vals, objectid)
				c++
			} else {
				cols += sep + config.DblQuote(key)
				p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
				sep = ","
				if key == joinField {
					j = strings.ToUpper(j.(string))

					if len(j.(string)) == 36 {
						j = "{" + j.(string) + "}"
					}

					//globalId = j.(string)
					//j = strings.Replace(j.(string), "}", "", -1)
					//j = strings.Replace(j.(string), "{", "", -1)
				}
				switch j.(type) {
				case float64:
					tmpFlt := j.(float64)
					if tmpFlt == float64(int(tmpFlt)) {
						vals = append(vals, int(tmpFlt))
					} else {
						vals = append(vals, j)
					}
				default:
					vals = append(vals, j)
				}
				c++
			}
		}
		if len(globalIdName) > 0 {
			cols += sep + config.DblQuote(globalIdName)
			p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
			vals = append(vals, uuidstr)
			i.Attributes[globalIdName] = uuidstr
			c++
		}
		//if config.Collector.Projects[name].Layers[id]["editFieldsInfo"] != nil {
		//if config.Collector.Projects[name].Layers[id].EditFieldsInfo != nil {
		//joinField = config.Collector.Projects[name].Layers[id]["joinField"].(string)
		current_time := time.Now().Local()

		//if rec, ok := config.Collector.Projects[name].Layers[id]["editFieldsInfo"].(map[string]interface{}); ok {
		if config.Collector.Projects[name].Layers[id].EditFieldsInfo != nil {
			//for key, j := range rec {
			//for key, j := range config.Collector.Projects[name].Layers[id].EditFieldsInfo {
			//for key, j := range config.Collector.Projects[name].Layers[id]["editFieldsInfo"] {
			//config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreatorField

			cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreatorField) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
			vals = append(vals, config.Collector.Username)
			p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
			i.Attributes["creatorField"] = config.Collector.Username
			c++

			cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].EditFieldsInfo.EditorField) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
			vals = append(vals, config.Collector.Username)
			p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
			i.Attributes["editorField"] = config.Collector.Username
			c++

			cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreationDateField)
			p += sep + config.Collector.DbTimeStamp                        //julianday('now')"
			i.Attributes["creationDateField"] = current_time.Unix() * 1000 //DateToNumber(current_time.Year(), current_time.Month(), current_time.Day())

			cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].EditFieldsInfo.EditDateField)
			p += sep + config.Collector.DbTimeStamp                    //julianday('now')"
			i.Attributes["editDateField"] = current_time.Unix() * 1000 //DateToNumber(current_time.Year(), current_time.Month(), current_time.Day())

			/*
				if key == "creatorField" || key == "editorField" {
					vals = append(vals, config.Collector.Username)
					p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
					i.Attributes[key] = config.Collector.Username
					c++
				} else if key == "creationDateField" || key == "editDateField" {
					p += sep + config.DbTimeStamp                  //julianday('now')"
					i.Attributes[key] = current_time.Unix() * 1000 //DateToNumber(current_time.Year(), current_time.Month(), current_time.Day())
					//year int, month time.Month, day int)
					//vals = append(vals, "julianday('now')")
					//cols += sep + j.(string) + "=julianday('now')"
				}
			*/

			//}
			//}

			/*
				cols += sep + config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creatorField"]
				p += sep + config.GetParam(c)
				c++

				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creatorField"] = config.Project.Username
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["editorField"]=config.Project.Username
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creationDateField"]=
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["editDateField"]
			*/
			/*
				if i.Geometry != nil {
					var geometry string
					cols += sep + config.DblQuote(config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreatorField) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
					vals = append(vals, geometry)
					p += sep + config.GetParam(config.Collector.DefaultDataSource, c)
					i.Attributes["creatorField"] = config.Collector.Username
					c++
				}
			*/
			//for key, j := range i.Geometry {
			//ST_Polygon('polygon ((52 28, 58 28, 58 23, 52 23, 52 28))', 4326)
			//ST_Point('point (52 24)', 4326)
			//}

		}

		//vals = append(vals, "")

		//cols += sep + joinField
		//p += sep + config.GetParam(c)
		//vals = append(vals, "")

		log.Println("insert into " + config.Collector.Schema + tableName + "(" + cols + ") values(" + p + ")")
		log.Print(vals)

		sql := "insert into " + config.Collector.Schema + tableName + "(" + cols + ") values(" + p + ")"
		/*
			stmt, err := config.GetReplicaDB(name).Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
		*/
		res, err := config.GetReplicaDB(name).Exec(sql, vals...)
		if err != nil {
			log.Println(err.Error())
		} else {
			if config.Collector.DefaultDataSource == structs.SQLITE3 {
				objectid, err := res.LastInsertId()
				if err != nil {
					println("Error:", err.Error())
				} else {
					println("LastInsertId:", objectid)
				}
			}
		}
		//stmt.Close()

		if config.Collector.DefaultDataSource == structs.PGSQL {
			//addsTxt = addsTxt[15 : len(addsTxt)-2]
			sql = "update " + config.Collector.Schema + "services set json=jsonb_set(json,'{features}',json->'features' || $1::jsonb,true) where type='query' and layerId=$2"
			log.Println(sql)
			stmt, err := config.Collector.DatabaseDB.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			//log.Println(i)
			//log.Println(id)
			var jsonstr []byte
			jsonstr, err = json.Marshal(i)
			if err != nil {
				log.Println(err)
			}

			_, err = stmt.Exec(jsonstr, id)
			if err != nil {
				log.Println(err.Error())
			}
			stmt.Close()
		} else if config.Collector.DefaultDataSource == structs.SQLITE3 {
			sql := "select json from " + config.Collector.Schema + "services where type='query' and layerId=?"
			stmt, err := config.Collector.Configuration.Prepare(sql)
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
			err = json.Unmarshal(row, &fieldObj)
			if err != nil {
				log.Println("Error unmarshalling fields into features object: " + string(row))
				log.Println(err.Error())
			}
			fieldObj.Features = append(fieldObj.Features, i)

			var jsonstr []byte
			jsonstr, err = json.Marshal(fieldObj)
			if err != nil {
				log.Println(err)
			}

			tx, err := config.Collector.Configuration.Begin()
			if err != nil {
				log.Fatal(err)
			}

			sql = "update " + config.Collector.Schema + "services set json=? where type='query' and layerId=?"

			stmt, err = tx.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}

			idInt, _ := strconv.Atoi(id)
			//log.Printf("JSON: %v:\n%v", string(jsonstr), idInt)

			_, err = tx.Stmt(stmt).Exec(string(jsonstr), idInt)
			if err != nil {
				log.Println(err.Error())
			}
			tx.Commit()
			stmt.Close()
		}
		result := map[string]interface{}{}
		result["objectId"] = objectid
		result["success"] = true
		result["globalId"] = nil

		results = append(results, result)
		objectid++
	}
	response, _ := json.Marshal(map[string]interface{}{"addResults": results, "updateResults": []string{}, "deleteResults": []string{}})
	return response
}

func AddsFile(name string, id string, parentTableName string, addsTxt string, joinField string, globalIdName string, parentObjectID string) []byte {
	current_time := time.Now().Local()
	jsonFile := config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json"
	//log.Println(jsonFile)
	file, err1 := ioutil.ReadFile(jsonFile)
	if err1 != nil {
		log.Println(err1)
	}
	var objectid int
	//var globalID string
	var results []interface{}

	var fieldObj structs.FeatureTable
	//map[string]map[string]map[string]
	err := json.Unmarshal(file, &fieldObj)
	if err != nil {
		log.Println("Error unmarshalling fields into features object: " + string(file))
		log.Println(err.Error())
	}

	var adds []structs.Feature
	decoder := json.NewDecoder(strings.NewReader(addsTxt)) //r.Body
	err = decoder.Decode(&adds)
	if err != nil {
		panic(err)
	}
	objectid = len(fieldObj.Features) + 1
	for _, i := range adds {
		//i.Attributes["objectId"] = objectid
		i.Attributes[parentObjectID] = objectid
		//i.Attributes["globalId"]=strings.ToUpper(i.Attributes["globalId"])
		if len(i.Attributes[joinField].(string)) > 0 {
			//input := strings.ToUpper(i.Attributes[joinField].(string))
			//tmpStr := input[1 : len(input)-1]
			i.Attributes[joinField] = strings.ToUpper(i.Attributes[joinField].(string))
			i.Attributes[joinField] = strings.Replace(i.Attributes[joinField].(string), "{", "", -1)
			i.Attributes[joinField] = strings.Replace(i.Attributes[joinField].(string), "}", "", -1)
			//strings.ToUpper(i.Attributes[joinField].(string)).Replace("{", "").Replace("{", "")
		}
		//if edit, save username and timestamp
		if config.Collector.Projects[name].Layers[id].EditFieldsInfo != nil {
			i.Attributes[config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreatorField] = config.Collector.Username
			i.Attributes[config.Collector.Projects[name].Layers[id].EditFieldsInfo.EditorField] = config.Collector.Username
			i.Attributes[config.Collector.Projects[name].Layers[id].EditFieldsInfo.CreationDateField] = current_time.Unix() * 1000
			i.Attributes[config.Collector.Projects[name].Layers[id].EditFieldsInfo.EditDateField] = current_time.Unix() * 1000
		}

		fieldObj.Features = append(fieldObj.Features, i)
		//write json back to file
		result := map[string]interface{}{}
		result["objectId"] = objectid
		result["success"] = true
		result["globalId"] = nil
		results = append(results, result)
		objectid++
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

	response, _ := json.Marshal(map[string]interface{}{"addResults": results, "updateResults": []string{}, "deleteResults": []string{}})
	return response
}
