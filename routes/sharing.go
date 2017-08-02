package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	config "github.com/traderboy/collector-server/config"
	structs "github.com/traderboy/collector-server/structs"
)

func sharing_info(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/rest/info")
	response, _ := json.Marshal(map[string]interface{}{"owningSystemUrl": "http://" + config.Server,
		"authInfo": map[string]interface{}{"tokenServicesUrl": "https://" + config.Collector.Hostname + "/sharing/rest/generateToken", "isTokenBasedSecurity": true}})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func sharing(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing")
	response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.Collector.ArcGisVersion})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	//w.Write(response)
	//setHeaders(c)
	//fmt.Println(response)
	//w.Write(response)
}
func sharing_rest(w http.ResponseWriter, r *http.Request) {

	log.Println("/sharing/rest (" + r.Method + ")")
	response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.Collector.ArcGisVersion})
	//w.Write(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func sharing_generateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rest := vars["rest"]

	log.Println("/sharing/" + rest + "/generateToken")
	/*
		tok := esritoken{
			Token:   "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp",
			Expires: 1940173783033,
			SSL:     false,
		}

		response, _ := json.Marshal(tok)
	*/
	var expires int64 = 1440173783033
	response, _ := json.Marshal(map[string]interface{}{"token": "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp", "expires": expires, "ssl": false})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func sharing_generateToken1(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/generateToken (post)")
	//response, _ := json.Marshal([]string{"token": "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp", "expires": 1940173783033, "ssl": false}
	/*
		tok := esritoken{
			Token:   "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp",
			Expires: 1940173783033,
			SSL:     false,
		}

		response, _ := json.Marshal(tok)
	*/
	var expires int64 = 1440173783033
	response, _ := json.Marshal(map[string]interface{}{"token": "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp", "expires": expires, "ssl": false})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
func sharing_authorize(w http.ResponseWriter, r *http.Request) {
	log.Println("//sharing/oauth2/authorize")
	log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "oauth2.html")
	http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"oauth2.html")
}
func sharing_authorize_v1(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/oauth2/authorize")
	http.Redirect(w, r, "/sharing/rest?f=json&culture=en-US&code=KIV31WkDhY6XIWXmWAc6U", http.StatusMovedPermanently)
	//302
	//c.Redirect(http.StatusMovedPermanently, "/sharing/rest?f=json&culture=en-US&code=KIV31WkDhY6XIWXmWAc6U")
	//http.ServeFile(w, r, config.Collector.DataPath + "/oauth2.html");
}
func sharing_approval(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/oauth2/approval")
	/*
		tok := esritoken{
			Token:   "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp",
			Expires: 1940173783033,
			SSL:     false,
		}
		response, _ := json.Marshal(tok)
	*/

	var expires int64 = 1440173783033
	response, _ := json.Marshal(map[string]interface{}{"token": "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp", "expires": expires, "ssl": false})

	w.Header().Set("Content-Type", "application/json")

	w.Write(response)
	//w.Write( response)
}
func sharing_signin(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/oauth2/signin")
	log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "search.json")
	http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"search.json")
}
func sharing_token(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/oauth2/token")

	var expires int64 = 99800
	response, _ := json.Marshal(map[string]interface{}{"access_token": config.AccessToken, "expires_in": expires, "username": "gisuser", "refresh_token": config.RefreshToken})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func sharing_tokens(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/rest/tokens")
	response, _ := json.Marshal(map[string]interface{}{"token": "1.0"})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func sharing_accounts_self(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, config.Collector.DataPath + "/search.json")
	log.Println("/sharing/{rest}/accounts/self (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcCatalog(body, "portals", "self", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}
	response := config.GetArcCatalog("portals", "self", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "portals.self.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"portals.self.json")
	}
}

/*
func sharing_accounts_self(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing//accounts/self (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcCatalog(body, "account", "self", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcCatalog("account", "self", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "account.self.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"account.self.json")
	}
}
*/
func sharing_portals_self(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/rest/portals/self (" + r.Method + ")")
	if r.Method == "PUT" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcCatalog(body, "portals", "self", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcCatalog("portals", "self", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "portals.self.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"portals.self.json")
	}
	//http.ServeFile(w, r, config.Collector.DataPath + "/portals_self.json")
}
func sharing_content_users(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	user := vars["user"]

	log.Println("/sharing/rest/content/users/" + user)
	//response, _ := json.Marshal([]string{ "username"{user}"),"total":0,"start":1,"num":0,"nextStart":-1,"currentFolder":nil,"items":[],"folders":[] }
	//folders := make([]int64], 0)
	//folders := make([]string], 0)
	response, _ := json.Marshal(map[string]interface{}{"folders": []string{}})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func sharing_content_items(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	//temp
	//name = config.ServiceName
	log.Println("/sharing/rest/content/items/" + name + "(" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcService(body, name, "content", -1, "items", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	//load from db
	response := config.GetArcService(name, "content", -1, "items", "")
	if len(response) > 0 {
		//log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.json")
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+"services"+string(os.PathSeparator)+"content.items.json")
	}
}

func sharing_content_items_data(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	/*
		if config.Collector.DefaultDataSource != structs.FILE {
			name = "%"
		}
	*/
	//log.Println("Old name:  " + name)
	//name = config.ServiceName
	//log.Println("New name:  " + name)
	log.Println("/sharing/rest/content/items/" + name + "/data (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcService(body, name, "content", -1, "data", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcService(name, "content", -1, "data", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.data.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"content.items.data.json")
	}
}
func sharing_search(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/rest/search (" + r.Method + ")")
	//vars := mux.Vars(r)

	//q := vars["q"]
	//q := r.Queries("q")
	q := r.FormValue("q")
	if strings.Index(q, "typekeywords") == -1 {
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "community", "groups", "")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("community", "groups", "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else {
			log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "community.groups.json")
			http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"community.groups.json")
		}
	} else {
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "search", "", "")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("search", "", "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else {
			log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "search.json")
			http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"search.json")
		}
	}
}
func sharing_community_users_notifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	log.Println("/sharing/rest/community/users/" + user + "/notifications")
	response, _ := json.Marshal(map[string]interface{}{"notifications": []string{}})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func sharing_community_users(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/rest/community/users/ (" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcCatalog(body, "community", "users", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcCatalog("community", "users", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "community.users.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"community.users.json")
	}
}

/*
func sharing_community_users(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	log.Println("/sharing//community/users/" + user + "(" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcCatalog(body, "community", "users", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcCatalog("community", "users", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	} else {

		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "community.users.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"community.users.json")
	}
}
*/
func sharing_community_users_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	log.Println("/sharing/rest/community/users/" + user + "(" + r.Method + ")")
	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
		ret := config.SetArcCatalog(body, "community", "users", "")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]interface{}{"response": ret})
		w.Write(response)
		return
	}

	response := config.GetArcCatalog("community", "users", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	} else {

		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "community.users.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"community.users.json")
	}
}

func sharing_community_users_user_info(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	img := vars["img"]
	log.Println("/sharing/rest/community/users/" + user + "/info/" + img)

	var path = "photos/cat.jpg"
	log.Println("Sending: " + path)
	http.ServeFile(w, r, path)
	//var fs = require("fs')
	//var file = fs.readFileSync(path, "utf8")
	//res.end(file)

}
func sharing_community_groups(w http.ResponseWriter, r *http.Request) {
	log.Println("/sharing/rest/community/groups")

	response := config.GetArcCatalog("community", "groups", "")
	if len(response) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + "community.groups.json")
		http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+"community.groups.json")
	}
}
func sharing_content_items_info_thumbnail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	//name := id
	//if config.Collector.DefaultDataSource != structs.FILE {
	//	id = "%"
	//}
	//log.Println("Old name:  " + id)
	//id = config.ServiceName
	//log.Println("New name:  " + id)

	img := vars["img"]
	if len(img) > 0 {
		log.Println("/sharing/rest/content/items/" + id + "/info/thumbnail/" + img)
	} else {
		img = "ago_downloaded.png"
		log.Println("/sharing/rest/content/items/" + id + "/info/thumbnail/ago_downloaded.png")
	}
	log.Println("Sending: " + config.Collector.DataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "thumbnails" + string(os.PathSeparator) + id + ".png")
	http.ServeFile(w, r, config.Collector.DataPath+string(os.PathSeparator)+id+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"thumbnails"+string(os.PathSeparator)+id+".png")
}
func sharing_content_items_info_thumbnail_v1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	//name := id
	if config.Collector.DefaultDataSource != structs.FILE {
		id = "%"
	}
	//log.Println("Old name:  " + id)
	id = config.ServiceName
	//log.Println("New name:  " + id)

	log.Println("/sharing/rest/content/items/" + id + "/info/thumbnail/ago_downloaded.png")
	response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.Collector.ArcGisVersion})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
