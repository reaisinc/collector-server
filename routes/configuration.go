package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	config "github.com/traderboy/collector-server/config"
)

func configuration(w http.ResponseWriter, r *http.Request) {
	log.Println("/config (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcCatalog(body, "config", "", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		
		log.Println("Reloading configuration")
		config.Initialize()
		return
	}

	response := config.GetArcCatalog("config", "", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + "config.json")
		http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+"config.json")
	}
}
