package routes

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func StartGorillaMux() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/config", configuration).Methods("GET", "PUT")
	r.HandleFunc("/info", info).Methods("GET")
	r.HandleFunc("/offline", offline).Methods("GET", "PUT")
	r.HandleFunc("/offline/{format}/{table}/{field}/{queryField}/{value}", offline_load).Methods("GET", "PUT")
	r.HandleFunc("/db", db).Methods("GET")

	//r.StrictSlash(true)

	//Download certs
	r.HandleFunc("/cert/", cert).Methods("GET")

	/*
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("/")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Welcome"))
		}).Methods("GET", "OPTIONS")
	*/
	r.HandleFunc("/arcgis", arcgis).Methods("GET")
	r.HandleFunc("/xml", xml).Methods("GET", "POST")

	/*
	   Root responses
	*/
	r.HandleFunc("/sharing", sharing).Methods("GET")
	r.HandleFunc("/sharing/", sharing).Methods("GET")
	r.HandleFunc("/sharing/rest", sharing_rest).Methods("GET", "POST")
	r.HandleFunc("/sharing/rest/", sharing_rest).Methods("GET", "POST", "PUT")

	//authentication.  Uses phoney tokens
	r.HandleFunc("/sharing/{rest}/generateToken", sharing_generateToken).Methods("GET", "POST")
	r.HandleFunc("/sharing/generateToken", sharing_generateToken).Methods("GET", "POST")

	r.HandleFunc("/sharing/oauth2/authorize", sharing_authorize).Methods("GET")
	r.HandleFunc("//sharing/oauth2/authorize", sharing_authorize).Methods("GET")

	r.HandleFunc("/sharing/oauth2/authorize_v1", sharing_authorize).Methods("GET")
	r.HandleFunc("/sharing/oauth2/approval", sharing_approval).Methods("GET")
	r.HandleFunc("/sharing/oauth2/signin", sharing_signin).Methods("GET")
	r.HandleFunc("/sharing/oauth2/token", sharing_token).Methods("GET", "POST")
	r.HandleFunc("/sharing/oauth2/token", sharing_token).Methods("GET", "POST")

	r.HandleFunc("/sharing/rest/tokens", sharing_tokens).Methods("GET", "POST")
	/*
	   openssl req -x509 -nodes -days 365 -newkey rsa:1024 \
	       -keyout /etc/ssl/private/reais.key \
	       -out /etc/ssl/certs/reais.crt
	*/

	/*
	   End authentication
	*/

	r.HandleFunc("/sharing/{rest}/accounts/self", sharing_accounts_self).Methods("GET", "PUT")
	r.HandleFunc("/sharing/accounts/self", sharing_accounts_self).Methods("GET", "PUT")
	//no customization necesssary except for username
	r.HandleFunc("/sharing/rest/portals/self", sharing_portals_self).Methods("GET", "PUT")
	r.HandleFunc("/sharing/rest/content/users/{user}", sharing_content_users).Methods("GET")
	r.HandleFunc("/sharing/rest/content/items/{name}", sharing_content_items).Methods("GET", "PUT")
	r.HandleFunc("/sharing/rest/content/items/{name}/data", sharing_content_items_data).Methods("GET", "PUT")
	r.HandleFunc("/sharing/rest/search", sharing_search).Methods("GET", "PUT")
	r.HandleFunc("/sharing/rest/community/users/{user}/notifications", sharing_community_users_notifications).Methods("GET")

	r.HandleFunc("/sharing//community/users/{user}", sharing_community_users_user).Methods("GET", "PUT")
	r.HandleFunc("/sharing/rest/community/users/{user}", sharing_community_users_user).Methods("GET", "PUT")
	r.HandleFunc("/sharing/rest/community/users", sharing_community_users).Methods("GET", "PUT")
	r.HandleFunc("/sharing/rest/community/users/{user}/info/{img}", sharing_community_users_user_info).Methods("GET")
	r.HandleFunc("/sharing/rest/community/groups", sharing_community_groups).Methods("GET", "POST")
	r.HandleFunc("/sharing/rest/content/items/{id}/info/thumbnail/{img}", sharing_content_items_info_thumbnail).Methods("GET")

	r.HandleFunc("/sharing/rest/content/items/{id}/info/thumbnail/ago_downloaded.png", sharing_content_items_info_thumbnail).Methods("GET")
	r.HandleFunc("/sharing/rest/info", sharing_info).Methods("GET")
	r.HandleFunc("/arcgis/rest/info", arcgis_info).Methods("GET")

	//Database functions
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/jobs/replicas", job_replicas).Methods("GET", "POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/unRegisterReplica", unRegisterReplica).Methods("POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/replicas/", replicas).Methods("GET", "POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/createReplica", createReplica).Methods("POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/synchronizeReplica", synchronizeReplica).Methods("POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/jobs", jobs).Methods("GET")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/jobs/jobs/", jobs_jobs).Methods("POST")
	r.HandleFunc("/arcgis/rest/services", services).Methods("GET", "PUT")
	r.HandleFunc("/arcgis/rest/services/", services).Methods("GET", "PUT")
	r.HandleFunc("/arcgis/rest/services", services_post).Methods("POST")
	r.HandleFunc("/arcgis/services", services_arcgis).Methods("GET")
	r.HandleFunc("/llarcgis/services", services_arcgis).Methods("POST")
	r.HandleFunc("/arcgis/services", services_post).Methods("POST")
	r.HandleFunc("/arcgis/rest/services//services/{name}/FeatureServer/info/metadata", info_metadata).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}", services_name).Methods("GET", "PUT")
	r.HandleFunc("/arcgis/rest/services//{name}", services_name).Methods("GET", "PUT")
	r.HandleFunc("/rest/services/{name}/FeatureServer", rest_services_name).Methods("GET", "POST", "PUT")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/uploads/upload", uploads_upload).Methods("POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer", featureServer).Methods("GET", "POST", "PUT")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}", id).Methods("GET", "POST", "PUT")

	//Attachments
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/attachments", attachment).Methods("GET")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/attachments/{img}", attachments_img).Methods("GET")
	/**
	Add attachments.
	File:  save attachment to filesystem only
	DB:    save attachment to filesystem and database
	*/
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/addAttachment", addAttachment).Methods("GET", "POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/updateAttachment", updateAttachment).Methods("POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/deleteAttachments", deleteAttachment).Methods("POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/db/{id}", db_id).Methods("GET", "POST", "PUT")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/xml/{id}", xml_id).Methods("GET", "POST", "PUT")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/table/{id}", table_id).Methods("GET", "POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/query", query).Methods("GET", "POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/queryRelatedRecords", queryRelatedRecords).Methods("GET", "POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/applyEdits", applyEdits).Methods("GET", "POST")
	//put this last - serve static content
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))
	return r
}

// ToNumber converts a year, month, day into a Julian Day Number.
// Based on http://en.wikipedia.org/wiki/Julian_day#Calculation.
// Only valid for dates after the Julian Day epoch which is January 1, 4713 BCE.
func DateToNumber(year int, month time.Month, day int) (julianDay int) {
	if year <= 0 {
		year++
	}
	a := int(14-month) / 12
	y := year + 4800 - a
	m := int(month) + 12*a - 3
	julianDay = int(day) + (153*m+2)/5 + 365*y + y/4
	if year > 1582 || (year == 1582 && (month > time.October || (month == time.October && day >= 15))) {
		return julianDay - y/100 + y/400 - 32045
	} else {
		return julianDay - 32083
	}
}
