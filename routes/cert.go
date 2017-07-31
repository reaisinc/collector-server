package routes

import (
	"log"
	"net/http"
	"os"

	config "github.com/traderboy/collector-server/config"
)

func cert(w http.ResponseWriter, r *http.Request) {
	//res.sendFile("certs/server.crt", { root : __dirname})
	log.Println("Sending: " + config.Collector.DataPath + "certs/server.crt")
	http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"certs"+string(os.PathSeparator)+"server.crt")
}
