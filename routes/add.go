package routes

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/traderboy/arcrestgo/structs"
	config "github.com/traderboy/collector-server/config"
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
	sql := "select \"base_id\"," + config.UUID + " from " + config.Schema + "\"GDB_RowidGenerators\" where \"registration_id\" in ( SELECT \"registration_id\" FROM " + config.Schema + "\"GDB_TableRegistry\" where \"table_name\"='" + parentTableName + "')"
	//sql := "select max(" + parentObjectID + ")+1," + config.UUID + " from " + tableName
	log.Println(sql)
	rows, err := config.DbQuery.Query(sql)
	//defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&objectid, &uuidstr)
		if err != nil {
			//log.Fatal(err)
			objectid = 1
			uuidstr = strings.ToUpper(uuid.Formatter(uuid.NewV4(), uuid.FormatCanonicalCurly))
		}
	}
	rows.Close()
	sql = "update " + config.Schema + "\"GDB_RowidGenerators\" set \"base_id\"=" + (strconv.Itoa(objectid + 1)) + " where \"registration_id\" in ( SELECT \"registration_id\" FROM " + config.Schema + "\"GDB_TableRegistry\" where \"table_name\"='" + parentTableName + "')"
	log.Println(sql)
	_, err = config.DbQuery.Exec(sql)
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
				p += sep + config.getParam(config.Collector.Projects[name].DataSource, c)
				sep = ","
				vals = append(vals, objectid)
				c++
			} else {
				cols += sep + config.DblQuote(key)
				p += sep + config.getParam(config.Collector.Projects[name].DataSource, c)
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
			p += sep + config.getParam(config.Collector.Projects[name].DataSource, c)
			vals = append(vals, uuidstr)
			i.Attributes[globalIdName] = uuidstr
			c++
		}
		if config.Collector.Projects[name].Layers[id]["editFieldsInfo"] != nil {
			//joinField = config.Collector.Projects[name].Layers[id]["joinField"].(string)
			current_time := time.Now().Local()

			if rec, ok := config.Collector.Projects[name].Layers[id]["editFieldsInfo"].(map[string]interface{}); ok {
				for key, j := range rec {
					//for key, j := range config.Collector.Projects[name].Layers[id]["editFieldsInfo"] {
					cols += sep + config.DblQuote(j.(string)) //config.Collector.Projects[name].Layers[id]["editFieldsInfo"][key]
					if key == "creatorField" || key == "editorField" {
						vals = append(vals, config.Collector.Username)
						p += sep + config.getParam(config.Collector.Projects[name].DataSource, c)
						i.Attributes[key] = config.Collector.Username
						c++
					} else if key == "creationDateField" || key == "editDateField" {
						p += sep + config.DbTimeStamp                  //julianday('now')"
						i.Attributes[key] = current_time.Unix() * 1000 //DateToNumber(current_time.Year(), current_time.Month(), current_time.Day())
						//year int, month time.Month, day int)
						//vals = append(vals, "julianday('now')")
						//cols += sep + j.(string) + "=julianday('now')"
					}

				}
			}

			/*
				cols += sep + config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creatorField"]
				p += sep + config.GetParam(c)
				c++

				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creatorField"] = config.Project.Username
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["editorField"]=config.Project.Username
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["creationDateField"]=
				config.Collector.Projects[name].Layers[id]["editFieldsInfo"]["editDateField"]
			*/

		}

		//vals = append(vals, "")

		//cols += sep + joinField
		//p += sep + config.GetParam(c)
		//vals = append(vals, "")

		log.Println("insert into " + config.Schema + tableName + "(" + cols + ") values(" + p + ")")
		log.Print(vals)

		sql := "insert into " + config.Schema + tableName + "(" + cols + ") values(" + p + ")"
		/*
			stmt, err := config.DbQuery.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
		*/
		res, err := config.DbQuery.Exec(sql, vals...)
		if err != nil {
			log.Println(err.Error())
		} else {
			if config.DbSource == config.SQLITE3 {
				objectid, err := res.LastInsertId()
				if err != nil {
					println("Error:", err.Error())
				} else {
					println("LastInsertId:", objectid)
				}
			}
		}
		//stmt.Close()

		if config.DbSource == config.PGSQL {
			//addsTxt = addsTxt[15 : len(addsTxt)-2]
			sql = "update " + config.Schema + "services set json=jsonb_set(json,'{features}',json->'features' || $1::jsonb,true) where type='query' and layerId=$2"
			log.Println(sql)
			stmt, err := config.Db.Prepare(sql)
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
		} else if config.DbSource == config.SQLITE3 {
			sql := "select json from " + config.Schema + "services where type='query' and layerId=?"
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
