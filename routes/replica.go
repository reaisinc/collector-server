package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	config "github.com/traderboy/collector-server/config"
)

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
	var fileName = config.Collector.Projects[name].ReplicaPath // + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
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
