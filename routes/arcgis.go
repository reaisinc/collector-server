package routes

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	config "github.com/traderboy/collector-server/config"
	structs "github.com/traderboy/collector-server/structs"
	"github.com/twinj/uuid"
)

func arcgis(w http.ResponseWriter, r *http.Request) {
	log.Println("/arcgis")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome"))
}
func arcgis_info(w http.ResponseWriter, r *http.Request) {
	log.Println("/arcgis/rest/info")
	response, _ := json.Marshal(map[string]interface{}{"currentVersion": "10.3", "fullVersion": "10.3", "authInfo": map[string]interface{}{"isTokenBasedSecurity": false}})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func job_replicas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/job/replica")
	var submissionTime int64 = 1441201696150
	var lastUpdatedTime int64 = 1441201705967
	response, _ := json.Marshal(map[string]interface{}{
		"replicaName": "MyReplica", "replicaID": "58808194-921a-4f9f-ac97-5ffd403368a9", "submissionTime": submissionTime, "lastUpdatedTime": lastUpdatedTime,
		"status": "Completed", "resultUrl": "http://" + config.Collector.Hostname + "/arcgis/rest/services/" + name + "/FeatureServer/replicas/"})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func replicas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/replicas")
	var fileName = config.Collector.Projects[name].ReplicaPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
	log.Println("Sending: " + fileName)
	http.ServeFile(w, r, fileName) //, { root : __dirname})
}
func unRegisterReplica(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/unRegisterReplica")
	response, _ := json.Marshal(map[string]interface{}{"success": true})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func createReplica(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/createReplica (post)")
	response, _ := json.Marshal(map[string]interface{}{"statusUrl": "http://" + config.Collector.Hostname + "/arcgis/rest/services/" + name + "/FeatureServer/replicas"})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func synchronizeReplica(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var deltaFile = "delta.geodatabase"
	replicaId := vars["replicaID"]

	//check to see if delta.geodatabase has been uploaded, then merge results with masterdatabase, then delete delta.geodatabase
	if _, err := os.Stat(deltaFile); !os.IsNotExist(err) {
		// path/to/whatever exists
		log.Println("Found delta.geodatabase")
		deltaDb, err := sql.Open("sqlite3", deltaFile)
		if err != nil {
			log.Fatal(err)
		}
		/*
		   SELECT a.*,b.*,'T_'||b.ChangedDatasetID||(case when b.ChangeType=2 then '_updates' else '_inserts' end)
		   FROM "GDB_DataChangesDatasets" a,"GDB_DataChangesDeltas" b
		   where a.id=b.id
		*/
		//sql := "ATTACH DATABASE '" + deltaFile + "' AS delta"

		sql := "SELECT \"ID\",\"LayerID\" FROM " + config.Collector.Schema + config.DblQuote("GDB_DataChangesDatasets")

		stmt, err := deltaDb.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
		}

		_, err = stmt.Exec()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
			log.Println(err.Error())
			return
		}
		stmt.Close()

		//sql := "DETACH DATABASE delta"
		/*
			err := os.Remove(deltaFile)
			if err != nil {
				log.Println("Unable to delete: " + deltaFile)
			}
		*/
	}

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/synchronizeReplica")
	//response, _ := json.Marshal(map[string]interface{}{"status": "Completed", "transportType": "esriTransportTypeUrl"})
	response, _ := json.Marshal(map[string]interface{}{"statusUrl": "http://" + config.Collector.Hostname + "/arcgis/rest/services/" + name + "/FeatureServer/jobs/" + replicaId})

	/*
		  "responseType": <esriReplicaResponseTypeEdits | esriReplicaResponseTypeEditsAndData| esriReplicaResponseTypeNoEdits>,
		  "resultUrl": "<url>", //path to JSON (dataFormat=JSON) or a SQLite geodatabase (dataFormat=sqlite)
		  "submissionTime": "<T1>",  //Time since epoch in milliseconds
		  "lastUpdatedTime": "<T2>", //Time since epoch in milliseconds
		  "status": "<Pending | InProgress | Completed | Failed | ImportChanges | ExportChanges | ExportingData | ExportingSnapshot
			       | ExportAttachments | ImportAttachments | ProvisioningReplica | UnRegisteringReplica | CompletedWithErrors>"
	*/
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
func jobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/jobs")
	var submissionTime int64 = 1441201696150
	var lastUpdatedTime int64 = 1441201705967
	response, _ := json.Marshal(map[string]interface{}{"replicaName": "MyReplica", "replicaID": "58808194-921a-4f9f-ac97-5ffd403368a9", "submissionTime": submissionTime,
		"lastUpdatedTime": lastUpdatedTime, "status": "Completed", "resultUrl": "http://" + config.Collector.Hostname + "/arcgis/rest/services/" + name + "/FeatureServer/replicas/"})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func jobs_jobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	jobs := vars["jobs"]

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/jobs/jobs")
	var submissionTime int64 = 1441201696150
	var lastUpdatedTime int64 = 1441201705967
	response, _ := json.Marshal(map[string]interface{}{"replicaName": "MyReplica", "replicaID": jobs, "submissionTime": submissionTime,
		"lastUpdatedTime": lastUpdatedTime, "status": "Completed", "resultUrl": ""})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func services(w http.ResponseWriter, r *http.Request) {
	log.Println("/arcgis/rest/services (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}

		ret := config.SetArcCatalog(body, "FeatureServer", "", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcCatalog("FeatureServer", "", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + "FeatureServer.json")
		http.ServeFile(w, r, config.Collector.DataPath+"FeatureServer.json")
	}
}

func services_arcgis(w http.ResponseWriter, r *http.Request) {
	log.Println("/arcgis/services")
	log.Println("Sending: " + config.Collector.DataPath + "FeatureServer.json")
	response := config.GetArcCatalog("FeatureServer", "", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		http.ServeFile(w, r, config.Collector.DataPath+"FeatureServer.json")
	}
}
func featureServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	dbPath := r.URL.Query().Get("db")

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcService(body, name, "FeatureServer", -1, "", dbPath)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcService(name, "FeatureServer", -1, "", dbPath)
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
	}
}

/*
func services_arcgis(w http.ResponseWriter, r *http.Request) {
	log.Println("/arcgis/services (post)")

	response := config.GetArcCatalog("FeatureServer", "", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + "FeatureServer.json")
		http.ServeFile(w, r, config.Collector.DataPath+"FeatureServer.json")
	}
}
*/
func services_post(w http.ResponseWriter, r *http.Request) {
	log.Println("/arcgis/rest/services (post)")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><head><title>Object moved</title></head><body>" +
		"<h2>Object moved to <a href=\"/arcgis/rest/services\">here</a>.</h2>" +
		"</body></html>"))
}
func info_metadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/info/metadata")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metadata stuff"))
}
func services_name(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	dbPath := r.URL.Query().Get("db")

	log.Println("/arcgis/rest/services/" + name + " (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcService(body, name, "", -1, "", dbPath)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcService(name, "FeatureServer", -1, "", dbPath)
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "FeatureServer.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"FeatureServer.json")
	}
}
func rest_services_name(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	dbPath := r.URL.Query().Get("db")

	log.Println("/rest/services/" + name + "/FeatureServer (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcService(body, name, "FeatureServer", -1, "", dbPath)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcService(name, "FeatureServer", -1, "", dbPath)
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
	}
}
func id(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id, _ := vars["id"]
	dbPath := r.URL.Query().Get("db")

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + " (" + r.Method + ")")

	idInt, _ := strconv.Atoi(id)
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcService(body, name, "FeatureServer", idInt, "", dbPath)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcService(name, "FeatureServer", idInt, "", dbPath)
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".json")
	}
}
func uploads_upload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/uploads/upload (" + r.Method + ")")
	/*
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcService(body, name, "FeatureServer", -1, "", dbPath)
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}
	*/
	const MAX_MEMORY = 10 * 1024 * 1024
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	for key, value := range r.MultipartForm.Value {
		//fmt.Fprintf(w, "%s:%s ", key, value)
		log.Printf("%s:%s", key, value)
	}
	//files, _ := ioutil.ReadDir(uploadPath)
	//fid := len(files) + 1
	var buf []byte
	var fileName string
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			fileName = fileHeader.Filename
			//path := fmt.Sprintf("%s%s%v%s%s", uploadPath, string(os.PathSeparator), objectid, "@", fileHeader.Filename)
			path := fmt.Sprintf("%s", fileHeader.Filename)
			log.Println(path)
			buf, _ = ioutil.ReadAll(file)
			ioutil.WriteFile(path, buf, os.ModePerm)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	/*
		{
		    "success": true,
		    "item": {
		        "itemID": "ib740c7bb-e5d0-4156-9cea-12fa7d3a472c",
		        "itemName": "lake.tif",
		        "description": "Lake Tahoe",
		        "date": 1246060800000,
		        "committed": true
		    }
		}
	*/
	//item, _ := json.Marshal(map[string]interface{}{"itemID": "1", "itemName": fileName, "description": "description", "date": time.Now().Local().Unix() * 1000, "committed": true})
	response, _ := json.Marshal(map[string]interface{}{"success": true, "item": map[string]interface{}{"itemID": "1", "itemName": fileName, "description": "description", "date": time.Now().Local().Unix() * 1000, "committed": true}})
	w.Write(response)

	/*
		response := config.GetArcService(name, "FeatureServer", -1, "", dbPath)
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
			http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
		}
	*/
}

func attachments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	row := vars["row"]
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/attachments")
	//{"attachmentInfos":[{"id":5,"globalId":"xxxx","parentID":"47","name":"cat.jpg","contentType":"image/jpeg","size":5091}]}
	//if config.Collector.Projects[name].AttachmentsPath == nil {
	//	config.Collector.Projects[name].AttachmentsPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
	//}
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
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
			return
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
	//var response []byte
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	//var response={"attachmentInfos":infos}

	//w.Write([]byte(AttachmentPath))
	/*
			if(fs.existsSync(AttachmentPath)){
			   var files = fs.readdirSync(AttachmentPath)
			   var infos=[]
			   for(var i in files)
			     infos.push({"id":i,"contentType":"image/jpeg","name":files[i]})
			     //{"id"{row}"),"contentType":"image/jpeg","name"{row}")+".jpg"}
			   response, _ := json.Marshal([]string{"attachmentInfos":infos}
			}
			else
			  response, _ := json.Marshal([]string{"attachmentInfos":[]}
		  w.Write(response)
	*/
}
func attachments_img(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	img := vars["img"]
	//imgInt, _ := strconv.Atoi(img)
	//img = strconv.Itoa(imgInt - 1)

	row := vars["row"]
	//imgInt, _ := strconv.Atoi(img)
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/attachments/img")
	//useFileSystem := false
	//if useFileSystem {
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
				return
			}
		}
		//{ "id": 2, "contentType": "application/pdf", "size": 270133,"name": "Sales Deed"  }
		response, _ := json.Marshal(map[string]interface{}{"error": "File not found"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
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
			w.Write(response)
			return
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

	/*
		files, _ := ioutil.ReadDir("./")
		for _, f := range files {
			fmt.Println(f.Name())
		}
	*/
	/*
		if _, err := os.Stat(attachment); err == nil {
			http.ServeFile(w, r, attachment)
		} else {
			response, _ := json.Marshal(map[string]interface{}{"status": "Completed", "transportType": "esriTransportTypeUrl"})
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}
	*/

	/*
			if(fs.existsSync(attachment))
		    res.sendFile(attachment)
		  else
		  	res.sendJSON({"Error":"File not found"})
	*/
	/*
			var path="photos/cat.jpg"
			var fs = require("fs')
		  var file = fs.readFileSync(path, "utf8")
		  res.end(file)
	*/
}
func addAttachment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	//idInt, _ := strconv.Atoi(id)
	row := vars["row"]

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/addAttachment")
	// TODO: move and rename the file using req.files.path & .name)
	//res.send(console.dir(req.files))  // DEBUG: display available fields
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
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
			return
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
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)
			return
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
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		const MAX_MEMORY = 10 * 1024 * 1024
		if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
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
		//w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
	/*
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(uploadPath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	*/

	/*
	   {
	     "addAttachmentResult" : {
	       "objectId" : 1,
	       "globalId" : "c9163a7c-f72b-472b-b495-902fde08c0be",
	       "success" : true
	     }
	   }
	*/

	/*
		  var mkdirp = require("mkdirp')
		  if(!fs.existsSync(uploadPath)){
		      //fs.mkdir(uploadPath,function(e){
		      mkdirp.sync(uploadPath,function(e){
		          if(!e || (e && e.code === 'EEXIST')){
		              //do something with contents

		          } else {
		              //debug
		              log.Println(e)
		          }
		      })
		  }
		  var files = fs.readdirSync(uploadPath)
		  var id=files.length
		  var fstream
		  req.pipe(req.busboy)
		  req.busboy.on("file", function (fieldname, file, filename) {
		        log.Println("Uploading: " + filename)
		        var attachment = uploadPath + "/" + id + ".jpg"
		        fstream = fs.createWriteStream(AttachmentPath)
		        file.pipe(fstream)
		        fstream.on("close", function () {
		            //res.redirect("back')
			          response, _ := json.Marshal([]string{"addAttachmentResult":{"objectId"{id},"globalId":null,"success":true}}
		            w.Write(response)
		        })
		  })
	*/
	/*
		  fs.readFile(req.files.attachment.path, function (err, data) {
		    // ...
		    var newPath = uploadPath + "/" + row") + ".jpg";
		    fs.writeFile(newPath, data, function (err) {
		      //res.redirect("back");
			    response, _ := json.Marshal([]string{"addAttachmentResult":{"objectId"{row}"),"globalId":null,"success":true}}
		      w.Write(response)
		    })
		  })
	*/
}
func updateAttachment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	idInt, _ := strconv.Atoi(id)
	row := vars["row"]
	var aid = r.FormValue("attachmentId")
	//aidInt, _ := strconv.Atoi(aid)
	//aid = strconv.Itoa(aidInt - 1)

	//if config.Collector.DefaultDataSource == structs.FILE {
	var uploadPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/updateAttachment")
	const MAX_MEMORY = 10 * 1024 * 1024
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func deleteAttachments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	row := vars["row"]

	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/deleteAttachments")
	var aid = r.FormValue("attachmentIds")
	aidInt, _ := strconv.Atoi(aid)
	//aid = strconv.Itoa(aidInt - 1)

	//results := []string{"objectId": id, "globalId": nil, "success": true}
	//results := []string{aid}

	var AttachmentPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
	files, _ := ioutil.ReadDir(AttachmentPath)
	//i := 0
	for _, f := range files {
		name := f.Name()
		if name[0:len(aid+"@")] == aid+"@" {
			err := os.Remove(AttachmentPath + string(os.PathSeparator) + f.Name())
			if err != nil {
				response, _ := json.Marshal(map[string]interface{}{"deleteAttachmentResults": aidInt, "error": err.Error()})
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
				return
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
func query(w http.ResponseWriter, r *http.Request) {
	//if(req.query.outFields=='OBJECTID'){
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	idInt, _ := strconv.Atoi(id)
	dbPath := r.URL.Query().Get("db")
	where := r.FormValue("where")
	outFields := r.FormValue("outFields")
	returnIdsOnly := r.FormValue("returnIdsOnly")
	var parentObjectID = config.Collector.Projects[name].Layers[id].Oidname
	//returnGeometry := r.FormValue("returnGeometry")
	objectIds := r.FormValue("objectIds")
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query")
	//returnIdsOnly = true

	//log.Println(r.FormValue("returnGeometry"))
	//log.Println(r.FormValue("outFields"))
	//sql := "select "+outFields + " from " +
	where = ""

	if len(where) > 0 {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/where=" + where)
		//response := config.GetArcQuery(name, "FeatureServer", idInt, "query",objectIds,where)
		w.Header().Set("Content-Type", "application/json")
		//var response = []byte("{\"objectIdFieldName\":\"OBJECTID\",\"globalIdFieldName\":\"GlobalID\",\"geometryProperties\":{\"shapeAreaFieldName\":\"Shape__Area\",\"shapeLengthFieldName\":\"Shape__Length\",\"units\":\"esriMeters\"},\"features\":[]}")
		var response = []byte(`{"objectIdFieldName":"OBJECTID","globalIdFieldName":"GlobalID","geometryProperties":{"shapeLengthFieldName":"","units":"esriMeters"},"features":[]}`)
		w.Write(response)

	} else if returnIdsOnly == "true" {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/objectids")

		response := config.GetArcService(name, "FeatureServer", idInt, "objectids", dbPath)
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".objectids.json")
			http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".objectids.json")
		}
	} else if len(objectIds) > 0 {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/objectIds=" + objectIds)

		//only get the select objectIds
		//response := config.GetArcService(name, "FeatureServer", idInt, "query")
		response := config.GetArcQuery(name, "FeatureServer", idInt, "query", parentObjectID, objectIds, dbPath)

		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
			http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")

		}
		//if returnGeometry == "false" &&
	} else if strings.Index(outFields, parentObjectID) > -1 { //r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/outfields=" + outFields)

		response := config.GetArcService(name, "FeatureServer", idInt, "outfields", dbPath)
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".outfields.json")
			http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".outfields.json")
		}
	} else {
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query/else")

		response := config.GetArcService(name, "FeatureServer", idInt, "query", dbPath)
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
			http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")
		}
	}

	//http.ServeFile(w, r, config.Collector.DataPath + "/" + id  + "query.json")

}

func queryRelatedRecords(w http.ResponseWriter, r *http.Request) {
	/*
		if 1 == 1 {
			//arcgis fields, arcgis vals
			var s = "{\"fields\":[{\"name\":\"OBJECTID\",\"type\":\"esriFieldTypeOID\",\"alias\":\"OBJECTID\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"occupied\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Occupation\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"condition_of_homesite\",\"type\":\"esriFieldTypeString\",\"alias\":\"Condition of homesite\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"solar_power\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Uses solar power?\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"septic_system\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Has septic system?\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_corrals\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of corrals\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_sheds\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of sheds\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_abandoned_vehicles\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of abandoned vehicles\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"structures_outside_boundary\",\"type\":\"esriFieldTypeString\",\"alias\":\"Structures outside homesite boundary\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"nonlessee_homesite_occupant\",\"type\":\"esriFieldTypeString\",\"alias\":\"Non-lessee homesite occupant_\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"condition_of_area\",\"type\":\"esriFieldTypeString\",\"alias\":\"Condition of area\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"lessee_denied_inspection\",\"type\":\"esriFieldTypeString\",\"alias\":\"Lessee denied inspection\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"Comments\",\"type\":\"esriFieldTypeString\",\"alias\":\"Comments\",\"sqlType\":\"sqlTypeOther\",\"length\":8000,\"domain\":null,\"defaultValue\":null},{\"name\":\"GlobalGUID\",\"type\":\"esriFieldTypeGUID\",\"alias\":\"GlobalGUID\",\"sqlType\":\"sqlTypeOther\",\"length\":38,\"domain\":null,\"defaultValue\":null},{\"name\":\"created_user\",\"type\":\"esriFieldTypeString\",\"alias\":\"Created user\",\"sqlType\":\"sqlTypeOther\",\"length\":255,\"domain\":null,\"defaultValue\":null},{\"name\":\"created_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Created date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"last_edited_user\",\"type\":\"esriFieldTypeString\",\"alias\":\"Last edited user\",\"sqlType\":\"sqlTypeOther\",\"length\":255,\"domain\":null,\"defaultValue\":null},{\"name\":\"last_edited_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Last edited date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_name\",\"type\":\"esriFieldTypeString\",\"alias\":\"Reviewer name\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Reviewer date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_title\",\"type\":\"esriFieldTypeString\",\"alias\":\"Reviewer title\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"GlobalID\",\"type\":\"esriFieldTypeGlobalID\",\"alias\":\"GlobalID\",\"sqlType\":\"sqlTypeOther\",\"length\":38,\"domain\":null,\"defaultValue\":null}],\"relatedRecordGroups\":[{\"objectId\":47,\"relatedRecords\":[{\"attributes\":{\"OBJECTID\":6,\"occupied\":1,\"condition_of_homesite\":null,\"solar_power\":null,\"septic_system\":null,\"number_corrals\":null,\"number_sheds\":null,\"number_abandoned_vehicles\":null,\"structures_outside_boundary\":null,\"nonlessee_homesite_occupant\":null,\"condition_of_area\":null,\"lessee_denied_inspection\":null,\"Comments\":null,\"GlobalGUID\":\"f66536f3-3f53-4cb1-8816-c7c366a02c8c\",\"created_user\":\"hpluser8\",\"created_date\":1499434034798,\"last_edited_user\":\"hpluser8\",\"last_edited_date\":1499434034798,\"reviewer_name\":null,\"reviewer_date\":null,\"reviewer_title\":null,\"GlobalID\":\"776b6cad-9427-47a4-a4a7-e81b701ef48e\"}}]}]}"
			//local fields, arcgis vals
			s = "{\"fields\":[{\"domain\":null,\"name\":\"OBJECTID\",\"nullable\":false,\"defaultValue\":null,\"editable\":false,\"alias\":\"OBJECTID\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeOID\"},{\"domain\":null,\"name\":\"occupied\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Occupation\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeSmallInteger\",\"length\":2},{\"domain\":null,\"name\":\"condition_of_homesite\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Condition of homesite\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":50},{\"domain\":null,\"name\":\"solar_power\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Uses solar power?\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeSmallInteger\",\"length\":2},{\"domain\":null,\"name\":\"septic_system\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Has septic system?\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeSmallInteger\",\"length\":2},{\"domain\":null,\"name\":\"number_corrals\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Number of corrals\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeSmallInteger\",\"length\":2},{\"domain\":null,\"name\":\"number_sheds\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Number of sheds\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeSmallInteger\",\"length\":2},{\"domain\":null,\"name\":\"number_abandoned_vehicles\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Number of abandoned vehicles\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeSmallInteger\",\"length\":2},{\"domain\":null,\"name\":\"structures_outside_boundary\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Structures outside homesite boundary\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":50},{\"domain\":null,\"name\":\"nonlessee_homesite_occupant\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Non-lessee homesite occupant_\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":50},{\"domain\":null,\"name\":\"condition_of_area\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Condition of area\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":50},{\"domain\":null,\"name\":\"lessee_denied_inspection\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Lessee denied inspection\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":50},{\"domain\":null,\"name\":\"Comments\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Comments\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":8000},{\"domain\":null,\"name\":\"GlobalGUID\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"GlobalGUID\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeGUID\",\"length\":38},{\"domain\":null,\"name\":\"created_user\",\"nullable\":true,\"defaultValue\":null,\"editable\":false,\"alias\":\"Created user\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":255},{\"domain\":null,\"name\":\"created_date\",\"nullable\":true,\"defaultValue\":null,\"editable\":false,\"alias\":\"Created date\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeDate\",\"length\":8},{\"domain\":null,\"name\":\"last_edited_user\",\"nullable\":true,\"defaultValue\":null,\"editable\":false,\"alias\":\"Last edited user\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":255},{\"domain\":null,\"name\":\"last_edited_date\",\"nullable\":true,\"defaultValue\":null,\"editable\":false,\"alias\":\"Last edited date\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeDate\",\"length\":8},{\"domain\":null,\"name\":\"reviewer_name\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Reviewer name\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":50},{\"domain\":null,\"name\":\"reviewer_date\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Reviewer date\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeDate\",\"length\":8},{\"domain\":null,\"name\":\"reviewer_title\",\"nullable\":true,\"defaultValue\":null,\"editable\":true,\"alias\":\"Reviewer title\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeString\",\"length\":50},{\"domain\":null,\"name\":\"GlobalID\",\"nullable\":false,\"defaultValue\":null,\"editable\":false,\"alias\":\"GlobalID\",\"sqlType\":\"sqlTypeOther\",\"type\":\"esriFieldTypeGlobalID\",\"length\":38}],\"relatedRecordGroups\":[{\"objectId\":47,\"relatedRecords\":[{\"attributes\":{\"OBJECTID\":6,\"occupied\":1,\"condition_of_homesite\":null,\"solar_power\":null,\"septic_system\":null,\"number_corrals\":null,\"number_sheds\":null,\"number_abandoned_vehicles\":null,\"structures_outside_boundary\":null,\"nonlessee_homesite_occupant\":null,\"condition_of_area\":null,\"lessee_denied_inspection\":null,\"Comments\":null,\"GlobalGUID\":\"f66536f3-3f53-4cb1-8816-c7c366a02c8c\",\"created_user\":\"hpluser8\",\"created_date\":1499434034798,\"last_edited_user\":\"hpluser8\",\"last_edited_date\":1499434034798,\"reviewer_name\":null,\"reviewer_date\":null,\"reviewer_title\":null,\"GlobalID\":\"776b6cad-9427-47a4-a4a7-e81b701ef48e\"}}]}]}"
			//arcgis fields, local vals
			s = "{\"fields\":[{\"name\":\"OBJECTID\",\"type\":\"esriFieldTypeOID\",\"alias\":\"OBJECTID\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"occupied\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Occupation\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"condition_of_homesite\",\"type\":\"esriFieldTypeString\",\"alias\":\"Condition of homesite\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"solar_power\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Uses solar power?\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"septic_system\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Has septic system?\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_corrals\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of corrals\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_sheds\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of sheds\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_abandoned_vehicles\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of abandoned vehicles\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"structures_outside_boundary\",\"type\":\"esriFieldTypeString\",\"alias\":\"Structures outside homesite boundary\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"nonlessee_homesite_occupant\",\"type\":\"esriFieldTypeString\",\"alias\":\"Non-lessee homesite occupant_\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"condition_of_area\",\"type\":\"esriFieldTypeString\",\"alias\":\"Condition of area\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"lessee_denied_inspection\",\"type\":\"esriFieldTypeString\",\"alias\":\"Lessee denied inspection\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"Comments\",\"type\":\"esriFieldTypeString\",\"alias\":\"Comments\",\"sqlType\":\"sqlTypeOther\",\"length\":8000,\"domain\":null,\"defaultValue\":null},{\"name\":\"GlobalGUID\",\"type\":\"esriFieldTypeGUID\",\"alias\":\"GlobalGUID\",\"sqlType\":\"sqlTypeOther\",\"length\":38,\"domain\":null,\"defaultValue\":null},{\"name\":\"created_user\",\"type\":\"esriFieldTypeString\",\"alias\":\"Created user\",\"sqlType\":\"sqlTypeOther\",\"length\":255,\"domain\":null,\"defaultValue\":null},{\"name\":\"created_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Created date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"last_edited_user\",\"type\":\"esriFieldTypeString\",\"alias\":\"Last edited user\",\"sqlType\":\"sqlTypeOther\",\"length\":255,\"domain\":null,\"defaultValue\":null},{\"name\":\"last_edited_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Last edited date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_name\",\"type\":\"esriFieldTypeString\",\"alias\":\"Reviewer name\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Reviewer date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_title\",\"type\":\"esriFieldTypeString\",\"alias\":\"Reviewer title\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"GlobalID\",\"type\":\"esriFieldTypeGlobalID\",\"alias\":\"GlobalID\",\"sqlType\":\"sqlTypeOther\",\"length\":38,\"domain\":null,\"defaultValue\":null}],\"relatedRecordGroups\":[{\"objectId\":\"47\",\"relatedRecords\":[{\"attributes\":{\"Comments\":null,\"GlobalGUID\":\"f66536f3-3f53-4cb1-8816-c7c366a02c8c\",\"GlobalID\":\"5ba5b963-bc05-4487-99a0-86e8d6fc5e3a\",\"OBJECTID\":1,\"condition_of_area\":null,\"condition_of_homesite\":null,\"created_date\":1499433868634,\"created_user\":\"hpluser8\",\"last_edited_date\":1499433868634,\"last_edited_user\":\"shale\",\"lessee_denied_inspection\":null,\"nonlessee_homesite_occupant\":null,\"number_abandoned_vehicles\":null,\"number_corrals\":null,\"number_sheds\":null,\"occupied\":1,\"reviewer_date\":null,\"reviewer_name\":null,\"reviewer_title\":null,\"septic_system\":null,\"solar_power\":null,\"structures_outside_boundary\":null}}]}]}"

			s = "{\"fields\":[{\"name\":\"OBJECTID\",\"type\":\"esriFieldTypeOID\",\"alias\":\"OBJECTID\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"occupied\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Occupation\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"condition_of_homesite\",\"type\":\"esriFieldTypeString\",\"alias\":\"Condition of homesite\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"solar_power\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Uses solar power?\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"septic_system\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Has septic system?\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_corrals\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of corrals\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_sheds\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of sheds\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_abandoned_vehicles\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of abandoned vehicles\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"structures_outside_boundary\",\"type\":\"esriFieldTypeString\",\"alias\":\"Structures outside homesite boundary\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"nonlessee_homesite_occupant\",\"type\":\"esriFieldTypeString\",\"alias\":\"Non-lessee homesite occupant_\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"condition_of_area\",\"type\":\"esriFieldTypeString\",\"alias\":\"Condition of area\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"lessee_denied_inspection\",\"type\":\"esriFieldTypeString\",\"alias\":\"Lessee denied inspection\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"Comments\",\"type\":\"esriFieldTypeString\",\"alias\":\"Comments\",\"sqlType\":\"sqlTypeOther\",\"length\":8000,\"domain\":null,\"defaultValue\":null},{\"name\":\"GlobalGUID\",\"type\":\"esriFieldTypeGUID\",\"alias\":\"GlobalGUID\",\"sqlType\":\"sqlTypeOther\",\"length\":38,\"domain\":null,\"defaultValue\":null},{\"name\":\"created_user\",\"type\":\"esriFieldTypeString\",\"alias\":\"Created user\",\"sqlType\":\"sqlTypeOther\",\"length\":255,\"domain\":null,\"defaultValue\":null},{\"name\":\"created_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Created date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"last_edited_user\",\"type\":\"esriFieldTypeString\",\"alias\":\"Last edited user\",\"sqlType\":\"sqlTypeOther\",\"length\":255,\"domain\":null,\"defaultValue\":null},{\"name\":\"last_edited_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Last edited date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_name\",\"type\":\"esriFieldTypeString\",\"alias\":\"Reviewer name\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Reviewer date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_title\",\"type\":\"esriFieldTypeString\",\"alias\":\"Reviewer title\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"GlobalID\",\"type\":\"esriFieldTypeGlobalID\",\"alias\":\"GlobalID\",\"sqlType\":\"sqlTypeOther\",\"length\":38,\"domain\":null,\"defaultValue\":null}],\"relatedRecordGroups\":[{\"objectId\":47",\"relatedRecords\":[{\"attributes\":{\"Comments\":null,\"GlobalGUID\":\"F66536F3-3F53-4CB1-8816-C7C366A02C8C\",\"GlobalID\":\"5BA5B963-BC05-4487-99A0-86E8D6FC5E3A\",\"OBJECTID\":4,\"condition_of_area\":null,\"condition_of_homesite\":null,\"created_date\":1499433868634,\"created_user\":\"shale\",\"last_edited_date\":1499433868634,\"last_edited_user\":\"shale\",\"lessee_denied_inspection\":null,\"nonlessee_homesite_occupant\":null,\"number_abandoned_vehicles\":null,\"number_corrals\":null,\"number_sheds\":null,\"occupied\":1,\"reviewer_date\":null,\"reviewer_name\":null,\"reviewer_title\":null,\"septic_system\":null,\"solar_power\":null,\"structures_outside_boundary\":null,}}]}]}"			//s = "{\"fields\":[{\"name\":\"OBJECTID\",\"type\":\"esriFieldTypeOID\",\"alias\":\"OBJECTID\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"occupied\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Occupation\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"condition_of_homesite\",\"type\":\"esriFieldTypeString\",\"alias\":\"Condition of homesite\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"solar_power\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Uses solar power?\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"septic_system\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Has septic system?\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_corrals\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of corrals\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_sheds\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of sheds\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"number_abandoned_vehicles\",\"type\":\"esriFieldTypeSmallInteger\",\"alias\":\"Number of abandoned vehicles\",\"sqlType\":\"sqlTypeOther\",\"domain\":null,\"defaultValue\":null},{\"name\":\"structures_outside_boundary\",\"type\":\"esriFieldTypeString\",\"alias\":\"Structures outside homesite boundary\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"nonlessee_homesite_occupant\",\"type\":\"esriFieldTypeString\",\"alias\":\"Non-lessee homesite occupant_\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"condition_of_area\",\"type\":\"esriFieldTypeString\",\"alias\":\"Condition of area\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"lessee_denied_inspection\",\"type\":\"esriFieldTypeString\",\"alias\":\"Lessee denied inspection\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"Comments\",\"type\":\"esriFieldTypeString\",\"alias\":\"Comments\",\"sqlType\":\"sqlTypeOther\",\"length\":8000,\"domain\":null,\"defaultValue\":null},{\"name\":\"GlobalGUID\",\"type\":\"esriFieldTypeGUID\",\"alias\":\"GlobalGUID\",\"sqlType\":\"sqlTypeOther\",\"length\":38,\"domain\":null,\"defaultValue\":null},{\"name\":\"created_user\",\"type\":\"esriFieldTypeString\",\"alias\":\"Created user\",\"sqlType\":\"sqlTypeOther\",\"length\":255,\"domain\":null,\"defaultValue\":null},{\"name\":\"created_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Created date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"last_edited_user\",\"type\":\"esriFieldTypeString\",\"alias\":\"Last edited user\",\"sqlType\":\"sqlTypeOther\",\"length\":255,\"domain\":null,\"defaultValue\":null},{\"name\":\"last_edited_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Last edited date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_name\",\"type\":\"esriFieldTypeString\",\"alias\":\"Reviewer name\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_date\",\"type\":\"esriFieldTypeDate\",\"alias\":\"Reviewer date\",\"sqlType\":\"sqlTypeOther\",\"length\":8,\"domain\":null,\"defaultValue\":null},{\"name\":\"reviewer_title\",\"type\":\"esriFieldTypeString\",\"alias\":\"Reviewer title\",\"sqlType\":\"sqlTypeOther\",\"length\":50,\"domain\":null,\"defaultValue\":null},{\"name\":\"GlobalID\",\"type\":\"esriFieldTypeGlobalID\",\"alias\":\"GlobalID\",\"sqlType\":\"sqlTypeOther\",\"length\":38,\"domain\":null,\"defaultValue\":null}],\"relatedRecordGroups\":[{\"objectId\":47,\"relatedRecords\":[{\"attributes\":{\"OBJECTID\":6,\"occupied\":1,\"condition_of_homesite\":null,\"solar_power\":null,\"septic_system\":null,\"number_corrals\":null,\"number_sheds\":null,\"number_abandoned_vehicles\":null,\"structures_outside_boundary\":null,\"nonlessee_homesite_occupant\":null,\"condition_of_area\":null,\"lessee_denied_inspection\":null,\"Comments\":null,\"GlobalGUID\":\"f66536f3-3f53-4cb1-8816-c7c366a02c8c\",\"created_user\":\"hpluser8\",\"created_date\":1499434034798,\"last_edited_user\":\"hpluser8\",\"last_edited_date\":1499434034798,\"reviewer_name\":null,\"reviewer_date\":null,\"reviewer_title\":null,\"GlobalID\":\"776b6cad-9427-47a4-a4a7-e81b701ef48e\"}}]}]}"
			//\"relatedRecordGroups\":[{\"objectId\":\"47\",\"relatedRecords\":[{\"attributes\":{\"Comments\":null,\"GlobalGUID\":\"f66536f3-3f53-4cb1-8816-c7c366a02c8c\",\"GlobalID\":\"776b6cad-9427-47a4-a4a7-e81b701ef48e\",\"OBJECTID\":4,\"condition_of_area\":null,\"condition_of_homesite\":null,\"created_date\":1499433868634,\"created_user\":\"shale\",\"last_edited_date\":1499433868634,\"last_edited_user\":\"shale\",\"lessee_denied_inspection\":null,\"nonlessee_homesite_occupant\":null,\"number_abandoned_vehicles\":null,\"number_corrals\":null,\"number_sheds\":null,\"occupied\":1,\"reviewer_date\":null,\"reviewer_name\":null,\"reviewer_title\":null,\"septic_system\":null,\"solar_power\":null,\"structures_outside_boundary\":null}}]}]}"
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(s))
			return
		}
	*/

	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	/*
		if id == "3" {
			jsonstr := `{"fields":[{"name":"OBJECTID","type":"esriFieldTypeOID","alias":"OBJECTID","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"cows","type":"esriFieldTypeSmallInteger","alias":"Cows","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"yearling_heifers","type":"esriFieldTypeSmallInteger","alias":"Yearling heifers","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"steer_calves","type":"esriFieldTypeSmallInteger","alias":"Steer calves","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"yearling_steers","type":"esriFieldTypeSmallInteger","alias":"Yearling steers","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"bulls","type":"esriFieldTypeSmallInteger","alias":"Bulls","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"mares","type":"esriFieldTypeSmallInteger","alias":"Mares","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"geldings","type":"esriFieldTypeSmallInteger","alias":"Geldings","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"studs","type":"esriFieldTypeSmallInteger","alias":"Studs","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"fillies","type":"esriFieldTypeSmallInteger","alias":"Fillies","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"colts","type":"esriFieldTypeSmallInteger","alias":"Colts","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"ewes","type":"esriFieldTypeSmallInteger","alias":"Ewes","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"lambs","type":"esriFieldTypeSmallInteger","alias":"Lambs","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"rams","type":"esriFieldTypeSmallInteger","alias":"Rams","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"wethers","type":"esriFieldTypeSmallInteger","alias":"Wethers","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"kids","type":"esriFieldTypeSmallInteger","alias":"Kids","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"billies","type":"esriFieldTypeSmallInteger","alias":"Billies","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"nannies","type":"esriFieldTypeSmallInteger","alias":"Nannies","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"Comments","type":"esriFieldTypeString","alias":"Comments","sqlType":"sqlTypeOther","length":8000,"domain":null,"defaultValue":null},{"name":"GlobalGUID","type":"esriFieldTypeGUID","alias":"GlobalGUID","sqlType":"sqlTypeOther","length":38,"domain":null,"defaultValue":null},{"name":"created_user","type":"esriFieldTypeString","alias":"Created user","sqlType":"sqlTypeOther","length":255,"domain":null,"defaultValue":null},{"name":"created_date","type":"esriFieldTypeDate","alias":"Created date","sqlType":"sqlTypeOther","length":8,"domain":null,"defaultValue":null},{"name":"last_edited_user","type":"esriFieldTypeString","alias":"Last edited user","sqlType":"sqlTypeOther","length":255,"domain":null,"defaultValue":null},{"name":"last_edited_date","type":"esriFieldTypeDate","alias":"Last edited date","sqlType":"sqlTypeOther","length":8,"domain":null,"defaultValue":null},{"name":"reviewer_name","type":"esriFieldTypeString","alias":"Reviewer name","sqlType":"sqlTypeOther","length":50,"domain":null,"defaultValue":null},{"name":"reviewer_date","type":"esriFieldTypeDate","alias":"Reviewer date","sqlType":"sqlTypeOther","length":8,"domain":null,"defaultValue":null},{"name":"reviewer_title","type":"esriFieldTypeString","alias":"Reviewer title","sqlType":"sqlTypeOther","length":50,"domain":null,"defaultValue":null},{"name":"GlobalID","type":"esriFieldTypeGlobalID","alias":"GlobalID","sqlType":"sqlTypeOther","length":38,"domain":null,"defaultValue":null}],"relatedRecordGroups":[]}`
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(jsonstr))
			return
		}
	*/

	var relationshipId = r.FormValue("relationshipId")
	var objectIds = r.FormValue("objectIds")
	var outFields = r.FormValue("outFields")
	var objectId, _ = strconv.Atoi(objectIds)
	//get fields for the related table
	dID := config.Collector.Projects[name].Relationships[relationshipId].DId
	var parentObjectID = config.Collector.Projects[name].Layers[id].Oidname

	//get the fields json

	if config.Collector.DefaultDataSource == structs.FILE {
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

		//var response []byte
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonstr)
		return
		//response = "{fields:" + fields + "," + response[1]
		//w.Write([]byte("{\"fields\":"))
		//w.Write(fields)
		//w.Write([]byte(","))
		//w.Write(response)
	}
	//idInt, _ := strconv.Atoi(id)

	var sql string
	var fields []byte
	var fieldsArr []structs.Field

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
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return
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
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return
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
	var response []byte
	if len(final_result) > 0 {
		var result = map[string]interface{}{}
		//result["objectId"] = objectIds //strconv.Atoi(objectIds)
		//OBS! must convert objectID to int or it fails on Android
		oid, _ := strconv.Atoi(objectIds)
		result["objectId"] = oid
		result["relatedRecords"] = final_result
		response, _ = json.Marshal(map[string]interface{}{"relatedRecordGroups": []interface{}{result}})
		response = response[1:]
	} else {
		response = []byte("\"relatedRecordGroups\":[]}")
	}
	//convert fields to string
	fields, err = json.Marshal(fieldsArr)
	if err != nil {
		log.Println(err)
	}

	//var response []byte
	w.Header().Set("Content-Type", "application/json")
	//response = "{fields:" + fields + "," + response[1]
	w.Write([]byte("{\"fields\":"))
	w.Write(fields)
	w.Write([]byte(","))
	w.Write(response)
	//w.Write([]byte("}"))
}

func applyEdits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	var parentObjectID = config.Collector.Projects[name].Layers[id].Oidname
	//idInt, _ := strconv.Atoi(id)
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/applyEdits")
	var response []byte
	var joinField = "GlobalID"
	//log.Println(config.Collector.Projects[name].Layers)
	//log.Println(config.Collector.Projects[name].Layers[id])
	//log.Println(config.Collector.Projects[name].Layers[id]["joinField"])
	if len(config.Collector.Projects[name].Layers[id].JoinField) > 0 {
		joinField = config.Collector.Projects[name].Layers[id].JoinField
	}
	if config.Collector.DefaultDataSource == structs.FILE {

		//get the fields json
		jsonFile := config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json"
		log.Println(jsonFile)
		file, err1 := ioutil.ReadFile(jsonFile)
		if err1 != nil {
			log.Println(err1)
		}
		var fieldObj structs.FeatureTable

		//map[string]map[string]map[string]
		err := json.Unmarshal(file, &fieldObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(file))
			log.Println(err.Error())
		}
		var objectid int
		//var globalID string
		var results []interface{}
		if len(r.FormValue("updates")) > 0 {
			var updates structs.Record
			decoder := json.NewDecoder(strings.NewReader(r.FormValue("updates"))) //r.Body
			err := decoder.Decode(&updates)
			if err != nil {
				panic(err)
			}
			//var objId int
			for k, i := range fieldObj.Features {
				//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
				//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))
				if int(i.Attributes[parentObjectID].(float64)) == objectid {
					//i.Attributes["OBJECTID"]
					fieldObj.Features[k].Attributes = updates[0].Attributes
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
			response, _ = json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": results, "deleteResults": []string{}})

			//response = Updates(name, id, tableName, r.FormValue("updates"))
		} else if len(r.FormValue("adds")) > 0 {
			//response = Adds(name, id, tableName, r.FormValue("adds"))
			var adds []structs.Feature
			decoder := json.NewDecoder(strings.NewReader(r.FormValue("adds"))) //r.Body
			err := decoder.Decode(&adds)
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

			response, _ = json.Marshal(map[string]interface{}{"addResults": results, "updateResults": []string{}, "deleteResults": []string{}})
		} else if len(r.FormValue("deletes")) > 0 {
			//response = Deletes(name, id, tableName, r.FormValue("deletes"))
			objectid, _ = strconv.Atoi(r.FormValue("deletes"))
			if objectid == 0 {
				return
			}
			for k, i := range fieldObj.Features {
				if int(i.Attributes[parentObjectID].(float64)) == objectid {
					//i.Attributes["OBJECTID"]
					fieldObj.Features = append(fieldObj.Features[:k], fieldObj.Features[k+1:]...)
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
			response, _ = json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": []string{}, "deleteResults": results})
		}

	} else {
		var tableName = config.Collector.Projects[name].Layers[id].Data
		var globalIdName = config.Collector.Projects[name].Layers[id].Globaloidname
		//log.Println("Table name: " + tableName)
		//var layerId = int(config.Services[name].Relationships[relationshipId]["dId"].(float64))

		if len(r.FormValue("updates")) > 0 {
			response = Updates(name, id, tableName, tableName+config.Collector.TableSuffix, r.FormValue("updates"), globalIdName, joinField, parentObjectID)
		} else if len(r.FormValue("adds")) > 0 {
			response = Adds(name, id, tableName, tableName+config.Collector.TableSuffix, r.FormValue("adds"), joinField, globalIdName, parentObjectID)
		} else if len(r.FormValue("deletes")) > 0 {
			response = Deletes(name, id, tableName, tableName+config.Collector.TableSuffix, r.FormValue("deletes"), globalIdName, parentObjectID)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	/*
		sql := "select json->'fields' from services where service=$1 and name=$2 and layerid=$3 and type=$4"
		log.Println(sql)
		log.Println("Values: " + name + "," + "FeatureServer" + "," + id)
		stmt, err := Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		var fields []byte
		err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
		if err != nil {
			log.Println(err.Error())
		}
	*/
	/*
		var replicaDb = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
		//var tableName = config.Services[name].Relationships[id]["dTable"].(string)
		//log.Println(tableName)
		//var layerId = int(config.Services[name].Relationships[id]["dId"].(float64))
		//id = "1"
		var jsonFile = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." +
			id + ".query.json"
		file, err1 := ioutil.ReadFile(jsonFile)
		if err1 != nil {
			fmt.Printf("// error while reading file %s\n", jsonFile)
			fmt.Printf("File error: %v\n", err1)
			os.Exit(1)
		}
	*/
	//var features map[string]interface{}{}
	//var features map[string]interface{}
	//var features map[string]map[string]map[string]map[string]interface{}
	//var features TableField
	/*
		var features []Field
		//map[string]map[string]map[string]
		err = json.Unmarshal(fields, &features)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		log.Println("Features dump:")
		log.Print(features)
		b, err1 := json.Marshal(features)
		if err1 != nil {
			log.Println(err1)
		}
		log.Println(string(b))
	*/

	//var replicaDb = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"

	//var jsonFields=JSON.parse(file)
	//log.Println("sqlite: " + replicaDb)
	//var db = new sqlite3.Database(replicaDb)
	/*
		var sqlstr = "select " + outFields + " from " +
			config.Services[name].Relationships[relationshipId]["dTable"].(string) +
			" where " +
			config.Services[name].Relationships[relationshipId]["dJoinKey"].(string) + " in (select " +
			config.Services[name].Relationships[relationshipId]["oJoinKey"].(string) + " from " +
			config.Services[name].Relationships[relationshipId]["oTable"].(string) +
			" where OBJECTID=$1)"
	*/

	/*
		var jsonOutputFile = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." +
			id + ".query.exported.json"

		os.Remove(jsonOutputFile)

		b, err1 := json.Marshal(features)
		if err1 != nil {
			log.Println(err1)
		}
		log.Println(string(b))
		ioutil.WriteFile(jsonOutputFile, b, 0644)
	*/
	//now read posted JSON
	//var updates = map[string]interface{}{}
}
