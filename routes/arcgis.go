package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	config "github.com/traderboy/collector-server/config"
	structs "github.com/traderboy/collector-server/structs"
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

func attachment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	row := vars["row"]
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/attachments")
	response := attachments(name, id, row)
	//{"attachmentInfos":[{"id":5,"globalId":"xxxx","parentID":"47","name":"cat.jpg","contentType":"image/jpeg","size":5091}]}
	//if config.Collector.Projects[name].AttachmentsPath == nil {
	//	config.Collector.Projects[name].AttachmentsPath = config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
	//}

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
	response := attachments_imgs(w, r, name, id, img, row)
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

	//useFileSystem := false
	//if useFileSystem {
	//w.Header().Set("Content-Type", contentType)
	//w.Write(response)

	//}

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
	response := addAttachments(r, name, id, row)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

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
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/updateAttachment")
	response := updateAttachments(r, name, id, idInt, row, aid)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func deleteAttachment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	row := vars["row"]
	var aid = r.FormValue("attachmentIds")
	aidInt, _ := strconv.Atoi(aid)
	log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/deleteAttachments")
	response := deleteAttachments(r, name, id, row, aid, aidInt)
	//aid = strconv.Itoa(aidInt - 1)
	//results := []string{"objectId": id, "globalId": nil, "success": true}
	//results := []string{aid}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func query(w http.ResponseWriter, r *http.Request) {
	//if(req.query.outFields=='OBJECTID'){
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	//idInt, _ := strconv.Atoi(id)
	//dbPath := r.URL.Query().Get("db")
	where := r.FormValue("where")
	outFields := r.FormValue("outFields")
	returnIdsOnly := r.FormValue("returnIdsOnly")
	var parentObjectID = config.Collector.Projects[name].Layers[id].Oidname
	var response []byte

	//returnGeometry := r.FormValue("returnGeometry")
	objectIds := r.FormValue("objectIds")

	/*
		<<<<<<< HEAD
				log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
				http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")
				return
		=======
			log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
			http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")
			return
		>>>>>>> f2ee24de79d7df3b1f9961b4452a18dfc07313b6
	*/

	//log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query")
	response = queryDB(name, "FeatureServer", id, where, outFields, returnIdsOnly, objectIds, parentObjectID)
	/*
		if config.Collector.DefaultDataSource != structs.FILE {
			//var response = []byte("{\"objectIdFieldName\":\"OBJECTID\",\"globalIdFieldName\":\"GlobalID\",\"geometryProperties\":{\"shapeAreaFieldName\":\"Shape__Area\",\"shapeLengthFieldName\":\"Shape__Length\",\"units\":\"esriMeters\"},\"features\":[]}")
			//var response = []byte(`{"objectIdFieldName":"OBJECTID","globalIdFieldName":"GlobalID","geometryProperties":{"shapeLengthFieldName":"","units":"esriMeters"},"features":[]}`)
			response = queryDB(name, "FeatureServer", id, where, outFields, returnIdsOnly, objectIds, parentObjectID)
			//log.Println(string(response))

		} else {
			where = ""
			response = queryText(name, "FeatureServer", id, where, outFields, returnIdsOnly, objectIds, parentObjectID)
		}
	*/
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	//returnIdsOnly = true

	//log.Println(r.FormValue("returnGeometry"))
	//log.Println(r.FormValue("outFields"))
	//sql := "select "+outFields + " from " +
	//where = ""
	/*
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
	*/
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
	var response []byte

	if config.Collector.DefaultDataSource == structs.FILE {
		response = queryRelatedRecordsFile(name, id, relationshipId, objectIds, objectId, outFields, parentObjectID, dID)

		//var response []byte

		//response = "{fields:" + fields + "," + response[1]
		//w.Write([]byte("{\"fields\":"))
		//w.Write(fields)
		//w.Write([]byte(","))
		//w.Write(response)
	} else {
		response = queryRelatedRecordsDB(name, id, relationshipId, objectIds, objectId, outFields, parentObjectID, dID)
		//var response []byte

	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	//idInt, _ := strconv.Atoi(id)

	//var response []byte
	/*
		w.Header().Set("Content-Type", "application/json")
		//response = "{fields:" + fields + "," + response[1]
		w.Write([]byte("{\"fields\":"))
		w.Write(fields)
		w.Write([]byte(","))
		w.Write(response)
	*/
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
	var tableName = config.Collector.Projects[name].Layers[id].Data
	var globalIdName = config.Collector.Projects[name].Layers[id].Globaloidname

	if config.Collector.DefaultDataSource == structs.FILE {

		//get the fields json
		/*
			current_time := time.Now().Local()
			jsonFile := config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json"
			//log.Println(jsonFile)
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
		*/
		if len(r.FormValue("updates")) > 0 {
			response = UpdatesFile(name, id, tableName, r.FormValue("updates"), globalIdName, joinField, parentObjectID)
		} else if len(r.FormValue("adds")) > 0 {
			response = AddsFile(name, id, tableName, r.FormValue("adds"), joinField, globalIdName, parentObjectID)
		} else if len(r.FormValue("deletes")) > 0 {
			response = DeletesFile(name, id, tableName, r.FormValue("deletes"), globalIdName, parentObjectID)
		}
		/*


			var objectid int
			//var globalID string
			var results []interface{}
			if len(r.FormValue("updates")) > 0 {
				UpdateFile()
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
		*/

	} else {
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

func runSqliteCmd(sql string, db string) string {
	exe := "sqlite3.exe"
	sql = "SELECT load_extension( 'stgeometry_sqlite.dll', 'SDE_SQL_funcs_init');" + sql + ";"
	args := []string{db, sql}
	var err error
	var out []byte
	out, err = exec.Command(exe, args...).Output()
	if err != nil {
		log.Println("Unable to execute sql command in sqlite3:  " + err.Error())
	}
	//if len(out) > 0 {
	//	log.Println(string(out))
	//}

	return strings.Trim(string(out), "\n\r")
}
