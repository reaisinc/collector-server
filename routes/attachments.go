package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	config "github.com/traderboy/collector-server/config"
	structs "github.com/traderboy/collector-server/structs"
	"github.com/twinj/uuid"
)

func addAttachments(r *http.Request, name string, id string, row string) []byte {
	var uploadPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
	os.MkdirAll(uploadPath, 0755)

	var objectid int
	var parentTableName = config.Collector.Projects[name].Layers[id].Data
	var parentObjectID = config.Collector.Projects[name].Layers[id].Oidname
	var tableName = parentTableName + "__ATTACH" + config.Collector.TableSuffix
	var globalIdName = config.Collector.Projects[name].Layers[id].Globaloidname
	var uuidstr string
	var globalid string
	log.Println("Table name: " + tableName)
	if config.Collector.DefaultDataSource == structs.FILE {
		var AttachmentPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
		files, _ := ioutil.ReadDir(AttachmentPath)
		//i := 0
		//find the largest ATTACHMENTID and inc
		objectid = 1
		globalid = strings.ToUpper(uuid.Formatter(uuid.NewV4(), uuid.FormatCanonicalCurly))
		globalid = globalid[1 : len(globalid)-1]
		for _, f := range files {
			name := f.Name()
			namearr := strings.Split(name, "@")

			if len(namearr) > 1 {
				curId, _ := strconv.Atoi(namearr[0])
				if curId > objectid {
					objectid = curId
				}
			}
			//if name[0:len(img+"@")] == img+"@" {
			//http.ServeFile(w, r, AttachmentPath+string(os.PathSeparator)+f.Name())
			//log.Println(AttachmentPath + string(os.PathSeparator) + f.Name())
			//return
			//}
		}

	} else {
		//sql := "select ifnull(max(ATTACHMENTID)+1,1) from " + tableName
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
		rows.Close()
		sql = "update " + config.Collector.Schema + "\"GDB_RowidGenerators\" set \"base_id\"=" + (strconv.Itoa(objectid + 1)) + " where \"registration_id\" in ( SELECT \"registration_id\" FROM " + config.Collector.Schema + "\"GDB_TableRegistry\" where \"table_name\"='" + parentTableName + "')"
		log.Println(sql)
		_, err = config.GetReplicaDB(name).Exec(sql)

		//log.Println(sql)
		//stmt, err := config.GetReplicaDB(name).Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			//w.Write([]byte(err.Error()))
			//w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			//w.Write(response)
			return response
		}
		//rows, err := config.GetReplicaDB(name).Query(sql)
		//err = stmt.QueryRow().Scan(&objectid)

		//get the parent globalid
		sql = "select " + config.DblQuote(globalIdName) + " from " + config.Collector.Schema + config.DblQuote(parentTableName) + " where " + config.DblQuote(parentObjectID) + "=" + config.GetParam(config.Collector.DefaultDataSource, 1)
		//log.Println(sql)
		stmt, err := config.GetReplicaDB(name).Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			//w.Write([]byte(err.Error()))
			//w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			//w.Write(response)
			return response
		}

		//rows, err := config.GetReplicaDB(name).Query(sql)
		err = stmt.QueryRow(row).Scan(&globalid)
		stmt.Close()
	} //END SQL
	/*
		cols += sep + key
		p += sep + config.GetParam(c)
		sep = ","
		vals = append(vals, objectid)
	*/

	//w.Write([]byte(uploadPath))
	/*
		if r.Method == "GET" {
			crutime := time.Now().Unix()
			h := md5.New()
			io.WriteString(h, strconv.FormatInt(crutime, 10))
			token := fmt.Sprintf("%x", h.Sum(nil))

			t, _ := template.ParseFiles("upload.gtpl")
			t.Execute(w, token)
		} else
		{
	*/
	const MAX_MEMORY = 10 * 1024 * 1024
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		log.Println(err)
		//http.Error(w, err.Error(), http.StatusForbidden)
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		return response
	}

	//for key, value := range r.MultipartForm.Value {
	//fmt.Fprintf(w, "%s:%s ", key, value)
	//log.Printf("%s:%s", key, value)
	//}
	//files, _ := ioutil.ReadDir(uploadPath)
	//fid := len(files) + 1
	var buf []byte
	var fileName string
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			fileName = fileHeader.Filename
			path := fmt.Sprintf("%s%s%v%s%s", uploadPath, string(os.PathSeparator), objectid, "@", fileHeader.Filename)
			log.Println(path)
			buf, _ = ioutil.ReadAll(file)
			ioutil.WriteFile(path, buf, os.ModePerm)
		}
	}
	if config.Collector.DefaultDataSource != structs.FILE {
		cols := "\"ATTACHMENTID\",\"GLOBALID\",\"REL_GLOBALID\",\"CONTENT_TYPE\",\"ATT_NAME\",\"DATA_SIZE\",\"DATA\"" //REL_GLOBALID
		sep := ""
		p := ""
		for i := 1; i < 8; i++ {
			p = p + sep + config.GetParam(config.Collector.DefaultDataSource, i)
			sep = ","
		}
		var vals []interface{}
		vals = append(vals, objectid)
		//vals = append(vals, config.Collector.UUID)
		vals = append(vals, uuidstr)
		vals = append(vals, globalid)
		vals = append(vals, http.DetectContentType(buf[:512]))
		vals = append(vals, fileName)
		vals = append(vals, len(buf))
		vals = append(vals, buf)

		//blob, err := ioutil.ReadAll(file)
		//c := 1

		//defer rows.Close()
		/*
			for rows.Next() {
				err := rows.Scan(&objectid)
				if err != nil {
					//log.Fatal(err)
					objectid = 1
				}
			}
			rows.Close()
		*/
		/*
			if len(globalIdName) > 0 {
				cols += sep + globalIdName
				p += sep + config.GetParam(c)
				vals = append(vals, globalId)
			}
		*/
		//1	{1085FDD1-89A3-4DEC-8171-787DA675FA84}	{89F39A8E-A4BD-4FB4-AE40-4A70F7AF6134}	image/jpeg	fark_EBoAgJdmC_knRWz-3t9Nx-2Tz8Y.jpg	21053	BLOB sz=21053 JPEG image
		//log.Println("insert into " + tableName + "(" + cols + ") values(" + p + ")")
		//log.Print(vals)

		sql := "insert into " + config.Collector.Schema + config.DblQuote(tableName) + "(" + cols + ") values(" + p + ")"
		log.Printf("insert into %v(%v) values(%v,'%v','%v','%v','%v',%v)", config.Collector.Schema+tableName, cols, vals[0], vals[1], vals[2], vals[3], vals[4], vals[5])

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

	}

	response, _ := json.Marshal(map[string]interface{}{"addAttachmentResult": map[string]interface{}{"objectId": objectid, "globalId": globalid, "success": true}})
	return response
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(response)
	//}
	//return nil
}

func updateAttachments(r *http.Request, name string, id string, idInt int, row string, aid string) []byte {
	var uploadPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
	const MAX_MEMORY = 10 * 1024 * 1024
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		log.Println(err)
		response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
		return response

		//http.Error(w, err.Error(), http.StatusForbidden)
	}

	for key, value := range r.MultipartForm.Value {
		log.Printf("%s:%s", key, value)
	}
	var buf []byte
	var fileName string
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			fileName = fileHeader.Filename
			path := fmt.Sprintf("%s%s%s%s%s", uploadPath, string(os.PathSeparator), aid, "@", fileHeader.Filename)
			log.Println(path)
			buf, _ = ioutil.ReadAll(file)
			ioutil.WriteFile(path, buf, os.ModePerm)
		}
	}
	//} else {
	if config.Collector.DefaultDataSource != structs.FILE {
		var parentTableName = config.Collector.Projects[name].Layers[id].Data
		var tableName = parentTableName + "__ATTACH" + config.Collector.TableSuffix

		cols := []string{"CONTENT_TYPE", "ATT_NAME", "DATA_SIZE", "DATA"}
		sep := ""
		p := ""
		for i := 0; i < len(cols); i++ {
			p = p + sep + config.DblQuote(cols[i]) + "=" + config.GetParam(config.Collector.DefaultDataSource, i)
			sep = ","
		}
		var vals []interface{}
		//vals = append(vals, objectid)
		//vals = append(vals, config.Collector.UUID)
		//vals = append(vals, globalid)

		vals = append(vals, http.DetectContentType(buf[:512]))
		vals = append(vals, fileName)
		vals = append(vals, len(buf))

		vals = append(vals, buf)

		sql := "update " + config.Collector.Schema + config.DblQuote(tableName) + " set " + p + " where " + config.DblQuote("ATTACHMENTID") + "=" + config.GetParam(config.Collector.DefaultDataSource, 1)
		log.Printf("update %v%v(%v) values('%v','%v',%v)", config.Collector.Schema, config.DblQuote(tableName), cols, vals[0], vals[1], vals[2])
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
	}

	/*
		var parentTableName = config.Collector.Schema + config.Collector.Projects[name].Layers[id].Data
		var tableName = parentTableName + "__ATTACH_evw"
		var vals []interface{}
		vals = append(vals, row)

		sql := "update " + tableName + " where OBJECTID=" + config.GetParam(0)
		log.Printf("delele from %v where OBJECTID=%v", tableName, row)
	*/

	//results[0] = gin.H{"objectId": id, "globalId": nil, "success": "true"}
	response, _ := json.Marshal(map[string]interface{}{"updateAttachmentResult": map[string]interface{}{"objectId": idInt, "globalId": nil, "success": true}})
	//w.Header().Set("Content-Type", "application/json")
	return response

}
func deleteAttachments(r *http.Request, name string, id string, row string, aid string, aidInt int) []byte {
	var AttachmentPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
	files, _ := ioutil.ReadDir(AttachmentPath)
	//i := 0
	for _, f := range files {
		name := f.Name()
		if name[0:len(aid+"@")] == aid+"@" {
			err := os.Remove(AttachmentPath + string(os.PathSeparator) + f.Name())
			if err != nil {
				response, _ := json.Marshal(map[string]interface{}{"deleteAttachmentResults": aidInt, "error": err.Error()})
				//w.Header().Set("Content-Type", "application/json")
				//w.Write(response)
				return response
			}
			log.Println("Deleting:  " + AttachmentPath + string(os.PathSeparator) + f.Name())
			break
		}
	}
	if config.Collector.DefaultDataSource != structs.FILE {
		var parentTableName = config.Collector.Projects[name].Layers[id].Data
		var parentObjectID = config.Collector.Projects[name].Layers[id].Oidname
		var tableName = parentTableName + "__ATTACH" + config.Collector.TableSuffix
		var vals []interface{}
		vals = append(vals, row)

		sql := "delete from " + config.Collector.Schema + config.DblQuote(tableName) + " where " + config.DblQuote("ATTACHMENTID") + "=" + config.GetParam(config.Collector.DefaultDataSource, 1)
		log.Printf("delele from %v where "+config.DblQuote(parentObjectID)+"=%v", tableName, row)

		_, err := config.GetReplicaDB(name).Exec(sql, vals...)
		if err != nil {
			log.Println(err.Error())
		}

	}

	response, _ := json.Marshal(map[string]interface{}{"deleteAttachmentResults": aidInt})
	return response
}
func attachments(name string, id string, row string) []byte {
	var AttachmentPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)

	//attachments:=[]interface{}
	attachments := make([]interface{}, 0)
	//[]interface{}
	//fields.Fields, "relatedRecordGroups": []interface{}{result}}
	//useFileSystem := false
	//if useFileSystem {
	if config.Collector.DefaultDataSource == structs.FILE {
		files, _ := ioutil.ReadDir(AttachmentPath)
		i := 0
		for _, f := range files {
			//tmpArr = strings.Split(f.Name(),"@")
			name = f.Name()
			idx := strings.Index(name, "@")
			if idx != -1 {
				fid, _ := strconv.Atoi(name[0:idx])
				//name = name[idx+1:]
				attachfile := map[string]interface{}{"id": fid, "contentType": "image/jpeg", "name": name[idx+1:]}
				attachments = append(attachments, attachfile)
			}
			i++
		}
	} else {
		//var objectid int
		//config.Collector.Schema +
		var parentTableName = config.Collector.Projects[name].Layers[id].Data
		var tableName = parentTableName + "__ATTACH" + config.Collector.TableSuffix
		var globalIdName = config.Collector.Projects[name].Layers[id].Globaloidname
		log.Println("Table name: " + tableName)

		sql := "select \"ATTACHMENTID\",\"CONTENT_TYPE\",\"ATT_NAME\" from " + config.Collector.Schema + config.DblQuote(tableName) + " where  " + config.DblQuote("REL_GLOBALID") + "=(select " + config.DblQuote(globalIdName) + " from " + config.Collector.Schema + config.DblQuote(parentTableName+config.Collector.TableSuffix) + " where " + config.DblQuote("OBJECTID") + "=" + config.GetParam(config.Collector.DefaultDataSource, 1) + ")"
		log.Printf("%v%v", sql, row)

		//stmt, err := config.GetReplicaDB(name).Prepare(sql)

		//rows, err := config.GetReplicaDB(name).Query(sql)
		var attachmentID int32
		var contentType string
		var attName string
		//err = stmt.QueryRow().Scan(&objectid)
		rows, err := config.GetReplicaDB(name).Query(sql, row)
		if err != nil {
			log.Println(err.Error())
			//w.Write([]byte(err.Error()))
			//w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			//w.Write(response)
			return response
		}

		for rows.Next() {
			err := rows.Scan(&attachmentID, &contentType, &attName)
			if err != nil {
				//log.Fatal(err)
				attachmentID = -1
			}
			attachfile := map[string]interface{}{"id": attachmentID, "contentType": contentType, "name": attName}
			attachments = append(attachments, attachfile)
		}
		rows.Close()
	}
	response, _ := json.Marshal(map[string]interface{}{"attachmentInfos": attachments})
	return response

}
func attachments_imgs(w http.ResponseWriter, r *http.Request, name string, id string, img string, row string) []byte {
	if config.Collector.DefaultDataSource == structs.FILE {

		//var attachment = config.AttachmentsPath + string(os.PathSeparator) + name + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator) + img + ".jpg"
		//var AttachmentPath = config.AttachmentsPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
		//if config.Collector.Projects[name].AttachmentsPath == nil {
		//	config.Collector.Projects[name].AttachmentsPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
		//}
		var AttachmentPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)

		files, _ := ioutil.ReadDir(AttachmentPath)
		//i := 0
		for _, f := range files {
			name := f.Name()
			if name[0:len(img+"@")] == img+"@" {
				http.ServeFile(w, r, AttachmentPath+string(os.PathSeparator)+f.Name())
				log.Println(AttachmentPath + string(os.PathSeparator) + f.Name())
				return []byte("")
			}
		}
		//{ "id": 2, "contentType": "application/pdf", "size": 270133,"name": "Sales Deed"  }
		response, _ := json.Marshal(map[string]interface{}{"error": "File not found"})
		return response
		//w.Header().Set("Content-Type", "application/json")
		//w.Write(response)
	} else {
		var parentTableName = config.Collector.Projects[name].Layers[id].Data
		var tableName = parentTableName + "__ATTACH" + config.Collector.TableSuffix
		var globalIdName = config.Collector.Projects[name].Layers[id].Globaloidname
		log.Println("Table name: " + tableName)

		sql := "select \"CONTENT_TYPE\",\"ATT_NAME\",\"DATA\" from " + config.Collector.Schema + config.DblQuote(tableName) + " where " + config.DblQuote("REL_GLOBALID") + "=(select " + config.DblQuote(globalIdName) + " from " + config.Collector.Schema + config.DblQuote(parentTableName+config.Collector.TableSuffix) + " where " + config.DblQuote("OBJECTID") + "=" + config.GetParam(config.Collector.DefaultDataSource, 1) + ")"
		log.Printf("%v%v", sql, row)

		//stmt, err := config.GetReplicaDB(name).Prepare(sql)

		//rows, err := config.GetReplicaDB(name).Query(sql)
		var attachment []byte
		var contentType string
		var attName string
		//err = stmt.QueryRow().Scan(&objectid)
		rows, err := config.GetReplicaDB(name).Query(sql, row)
		if err != nil {
			log.Println(err.Error())
			//w.Write([]byte(err.Error()))
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			//w.Write(response)
			return response
		}

		for rows.Next() {
			err := rows.Scan(&contentType, &attName, &attachment)
			if err != nil {
				//log.Fatal(err)

			}
			//attachfile := map[string]interface{}{"id": attachmentID, "contentType": contentType, "name": attName}
			//attachments = append(attachments, attachfile)
		}
		rows.Close()
		w.Header().Set("Content-Type", contentType)

		w.Write(attachment)
	}
	return []byte("")

	//return response
}
