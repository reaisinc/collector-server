package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	sqlite3 "github.com/mattn/go-sqlite3"
	structs "github.com/traderboy/collector-server/structs"
)

/*
const (
	PGSQL   = "pgsql"
	SQLITE3 = "sqlite"
	FILE    = "file"
)
*/
//var Catalogs map[string]structs.Catalog

//var DbSource = 0 // SQLITE3

//var Schema = "" //= "postgres."
//var TableSuffix = ""
//var DbTimeStamp = ""

var Collector structs.Collector

//var Project structs.Project
//var DataPath = "catalogs"
var SqlFlags = "?cache=shared&mode=wrc"
var SqlWalFlags = "?PRAGMA journal_mode=WAL"

//leasecompliance2016
var ServiceName string
var HTTPPort string  // = ":80"
var HTTPSPort string //= ":443"
var Pem string
var Cert string
var UUID = ""

//"github.com/gin-gonic/gin"
//Db is the SQLITE databa se object

//var configFile = DataPath + string(os.PathSeparator) + "config.json"
//var ArcGisVersion = "3.8"

//var Db *sql.DB
//var DbQuery *sql.DB
//var DbSqliteQuery *sql.DB
//var DbSqliteDbName string

//var port = ":8080"

//var DataPath = RootPath        //+ string(os.PathSeparator)        //+ string(os.PathSeparator) //+ "services"
//var ReplicaPath = DataPath     //+ string(os.PathSeparator)     //+ "replicas"
//var AttachmentsPath = DataPath //+ string(os.PathSeparator) //+ "attachments"

//var CertificatePath = "ssl" + string(os.PathSeparator) + "agent2-cert.cert"

//var config map[string]interface{}
//var defaultService = ""
//var UploadPath = ""
var Server = ""
var RefreshToken = "51vzPXXNl7scWXsw7YXvhMp_eyw_iQzifDIN23jNSsQuejcrDtLmf3IN5_bK0P5Z9K9J5dNb2yBbhXqjm9KlGtv5uDjr98fsUAAmNxGqnz3x0tvl355ZiuUUqqArXkBY-o6KaDtlDEncusGVM8wClk0bRr1-HeZJcR7ph9KU9khoX6H-DcFEZ4sRdl9c16exIX5lGIitw_vTmuomlivsGIQDq9thskbuaaTHMtP1m3VVnhuRQbyiZTLySjHDR8OVllSPc2Fpt0M-F5cPl_3nQg.."
var AccessToken = "XMdOaajM4srQWx8nQ77KuOYGO8GupnCoYALvXEnTj0V_ZXmEzhrcboHLb7hGtGxZCYUGFt07HKOTnkNLah8LflMDoWmKGr4No2LBSpoNkhJqc9zPa2gR3vfZp5L3yXigqxYOBVjveiuarUo2z_nqQ401_JL-mCRsXq9NO1DYrLw."
var once sync.Once

func Initialize() {
	var DataPath = "catalogs"
	SqlFlags = ""
	SqlWalFlags = ""
	//clear out variables in case Initialize is run again

	//DbName = ""
	Cert = ""
	Pem = ""
	HTTPPort = ""
	HTTPSPort = ""
	var DataSource = ""

	//var err error
	//pwd, err := os.Getwd()
	//if err != nil {
	//	log.Println("Unable to get current directory")
	//}
	//DataPath = pwd + string(os.PathSeparator) + DataPath //+ string(os.PathSeparator)
	//var err error
	//var DbName string
	//if len(os.Getenv("DB_SOURCE")) > 0 {
	//read in from environment variables
	if len(os.Getenv("DATA_PATH")) > 0 {
		DataPath, _ = filepath.Abs(os.Getenv("DATA_PATH"))
	}
	//set default values if missing in config/cli/env vars
	//}
	//for docker, environment variables override command line parameters?
	if len(os.Getenv("HTTP_PORT")) > 0 {
		HTTPPort = ":" + os.Getenv("HTTP_PORT") //80
	}

	if len(os.Getenv("HTTPS_PORT")) > 0 {
		HTTPSPort = ":" + os.Getenv("HTTPS_PORT") //443
	}

	if len(os.Getenv("PEM_PATH")) > 0 {
		Pem = os.Getenv("PEM_PATH")
	}
	if len(os.Getenv("CERT_PATH")) > 0 {
		Cert = os.Getenv("CERT_PATH")
	}

	//now override any settings from command line or environment variables
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			//log.Println(os.Args[i][0] == 45)
			/*
					if os.Args[i] == "-sqlite" {
						DbSource = SQLITE3
						if len(os.Args) > i+1 && os.Args[i+1][0] != 45 { //&& len(os.Args[i+1]) > 0 && os.Args[i+1][0] != 45
							DbName = os.Args[i+1]
						} else {
							//DbName = pwd + string(os.PathSeparator) + "catalogs/collectorDb.sqlite"
							fmt.Println("No sqlite path entered")
							fmt.Println("Example:  catalogs/collectorDb.sqlite")
							os.Exit(1)
						}
					} else if os.Args[i] == "-pgsql" {
						DbSource = PGSQL
						if len(os.Args) > i+1 && os.Args[i+1][0] != 45 { // && len(os.Args[i+1]) > 0 && os.Args[i+1][0] != 45
							DbName = os.Args[i+1]
						} else {
							fmt.Println("No Postgresql connection string entered")
							fmt.Println("Example:  user=postgres dbname=gis host=192.168.99.100")
							os.Exit(1)

							//DbName = "user=postgres dbname=gis host=192.168.99.100"
						}
					} else if os.Args[i] == "-data" {
						if len(os.Args) > i+1 && os.Args[i+1][0] != 45 {
							DataPath, _ = filepath.Abs(os.Args[i+1])
						} else {
							fmt.Println("No data path to catalogs entered")
							os.Exit(1)
						}
						//ServiceName = filepath.Base(os.Args[i+1])
					} else
				} else if os.Args[i] == "-file" {
					DbSource = FILE
					//LoadConfigurationFromFile()
			*/
			if os.Args[i] == "-file" {
				DataSource = structs.FILE
			}
			if os.Args[i] == "-p" && len(os.Args) > i+1 {
				HTTPPort = ":" + os.Args[i+1]
			} else if os.Args[i] == "-https" && len(os.Args) > i && len(os.Args[i+1]) > 0 {
				HTTPSPort = ":" + os.Args[i+1]
			} else if os.Args[i] == "-pem" && len(os.Args) > i && len(os.Args[i+1]) > 0 {
				Pem = os.Args[i+1]
			} else if os.Args[i] == "-cert" && len(os.Args) > i && len(os.Args[i+1]) > 0 {
				Cert = os.Args[i+1]
			} else if os.Args[i] == "-data" {
				if len(os.Args) > i+1 && os.Args[i+1][0] != 45 {
					DataPath, _ = filepath.Abs(os.Args[i+1])
				} else {
					fmt.Println("Invalid data path to catalogs entered: " + string(os.Args[i+1][0]))
					os.Exit(1)
				}
				//ServiceName = filepath.Base(os.Args[i+1])
			} else if os.Args[i] == "-h" {
				fmt.Println("Usage:")
				fmt.Println("go run server.go -p HTTP Port -https HTTPS Port -data <path to catalogs folder> -sqlite <path to service .sqlite> -pgsql <connection string for Postgresql> -pem <path to pem> -cert <path to cert> -h [show help]")
				os.Exit(0)
			}
		}
	}

	/*

		tmpSrc := os.Getenv("DB_SOURCE")
		if len(tmpSrc) > 0 {
			if tmpSrc == "PGSQL" {
				DbSource = PGSQL
				//DbName = os.Getenv("DB_NAME")
			} else if tmpSrc == "SQLITE" {
				DbSource = SQLITE3
				//DbName = os.Getenv("DB_NAME")
			} else {
				DbSource = FILE
			}
		} else {
			DbSource = FILE
		}
	*/
	/*
		} else if _, err := os.Stat("catalogs/collectorDb.sqlite"); !os.IsNotExist(err) {
			DbName = "catalogs/collectorDb.sqlite"

			Db, err = sql.Open("sqlite3", DbName+SqlFlags)
			if err != nil {
				log.Fatal(err)
			}

			//get configuration
			LoadConfiguration()
			//Db.Close()
			tmpSrc := Project.DataSource
			if len(tmpSrc) > 0 {
				if tmpSrc == "psql" {
					DbSource = PGSQL
					//DbName = os.Getenv("DB_NAME")
				} else if tmpSrc == "sqlite" {
					DbSource = SQLITE3
					//DbName = os.Getenv("DB_NAME")
				} else {
					DbSource = FILE
				}
			} else {
				DbSource = FILE
			}
			if len(HTTPPort) == 0 {
				HTTPPort = ":" + Project.HttpPort
				if len(HTTPPort) == 1 {
					HTTPPort = ":80"
				}
			}
			if len(HTTPSPort) == 0 {
				HTTPSPort = ":" + Project.HttpsPort
				if len(HTTPSPort) == 1 {
					HTTPSPort = ":443"
				}
			}
			if len(Pem) == 0 {
				Pem = Project.Pem
			}
			if len(Cert) == 0 {
				Cert = Project.Cert
			}
	*/
	//} //else {
	//DataPath, _ = filepath.Abs(DataPath)

	//RootPath, _ = filepath.Abs(os.Args[i+1])
	//ServiceName = filepath.Base(os.Args[i+1])
	/*
		if Project.DataSource == "pg" {
			DbSource = PGSQL
			DbName = Project.PG
			//Schema = "postgres."
		} else if Project.DataSource == "sqlite" {
			DbSource = SQLITE3
			DbName = Project.SqliteDb
		} else if Project.DataSource == "file" {
			DbSource = FILE
		}
	*/
	/*
		tmpSrc := Project.DataSource
		if len(tmpSrc) > 0 {
			if tmpSrc == "pgsql" {
				DbSource = PGSQL
				if len(DbName) == 0 {
					DbName = Project.PG
				}
				//DbName = os.Getenv("DB_NAME")
			} else if tmpSrc == "sqlite" {
				DbSource = SQLITE3
				if len(DbName) == 0 {
					DbName = Collector.SqliteDb
				}
				if len(DbName) == 0 {
					DbSource = FILE
				}
				//DbName = os.Getenv("DB_NAME")
			} else {
				DbSource = FILE
			}
		} else {
			DbSource = FILE
		}
	*/

	//Load the config.json file first, then override any
	DataPath, _ = filepath.Abs(DataPath)
	if _, err := os.Stat(DataPath + string(os.PathSeparator) + "config.json"); os.IsNotExist(err) {
		DataPath = DataPath + string(os.PathSeparator) + "catalogs"
		if _, err := os.Stat(DataPath + string(os.PathSeparator) + "config.json"); os.IsNotExist(err) {
			fmt.Println("Unable to locate catalogs directory in data path: " + DataPath)
			os.Exit(0)
		}
	}

	//read all folder in catalogs
	/*
		files, _ := ioutil.ReadDir(RootPath)

		for _, f := range files {
			if f.IsDir() {
		}
		}
	*/

	LoadConfigurationFromFile(DataPath)
	//Collector.DataPath = DataPath
	if len(DataSource) > 0 {
		Collector.DefaultDataSource = DataSource
	}

	//override any settings from config file
	if len(Collector.HttpPort) == 1 {
		Collector.HttpPort = ":80"
	} else {
		Collector.HttpPort = ":" + Collector.HttpPort
	}
	if len(Collector.HttpsPort) == 1 {
		Collector.HttpsPort = ":443"
	} else {
		Collector.HttpsPort = ":" + Collector.HttpsPort
	}

	if len(Pem) > 0 {
		Collector.Pem = Pem
	}
	if len(Cert) > 0 {
		Collector.Cert = Cert
	}
	//overwrite if using openshift 2
	if len(os.Getenv("OPENSHIFT_GO_IP")) > 0 {
		Collector.Hostname = os.Getenv("OPENSHIFT_GO_IP")
	}
	if len(os.Getenv("OPENSHIFT_GO_PORT")) > 0 {
		Collector.HttpPort = os.Getenv("OPENSHIFT_GO_PORT")
	}

	var err error
	if Collector.DefaultDataSource != structs.FILE {
		//connect to configuration sqlite database
		Collector.Configuration, err = sql.Open("sqlite3", Collector.SqliteDb+SqlFlags)
		if err != nil {
			log.Fatal(err)
		}
		err = Collector.Configuration.Ping()
		if err != nil {
			log.Fatalf("Error on opening database connection: %s", err.Error())
		}
	}

	if Collector.DefaultDataSource == structs.SQLITE3 {
		Collector.Schema = ""
		Collector.TableSuffix = "_evw"
		Collector.UUID = "(select '{'||upper(substr(u,1,8)||'-'||substr(u,9,4)||'-4'||substr(u,13,3)||'-'||v||substr(u,17,3)||'-'||substr(u,21,12))||'}' from ( select lower(hex(randomblob(16))) as u, substr('89ab',abs(random()) % 4 + 1, 1) as v) as foo)"
		Collector.DbTimeStamp = "(julianday('now') - 2440587.5)*86400.0*1000"
		Collector.DatabaseDB = Collector.Configuration

		//once.Do(initDB)
		/*
			Collector.DatabaseDB, err = sql.Open("sqlite3", Collector.SqliteDb+SqlFlags)
			if err != nil {
				log.Fatal(err)
			}
			err = Collector.DatabaseDB.Ping()
			if err != nil {
				log.Fatalf("Error on opening database connection: %s", err.Error())
			}
		*/
		//DbQueryName := DataPath + string(os.PathSeparator) + ServiceName + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + ServiceName + ".geodatabase"

		//DbQuery, err = sql.Open("sqlite3", "file:"+DbQueryName+"?PRAGMA journal_mode=WAL")
		/*
			Collector.Replica, err = sql.Open("sqlite3", DbQueryName+SqlFlags)
			if err != nil {
				log.Fatal(err)
			}
			err = DbQuery.Ping()
			if err != nil {
				log.Fatalf("Error on opening database connection: %s", err.Error())
			}
			log.Println("Sqlite replica database: " + DbQueryName)
		*/

	} else if Collector.DefaultDataSource == structs.PGSQL {
		Collector.DatabaseDB, err = sql.Open("postgres", Collector.PG)
		if err != nil {
			log.Fatal(err)
		}
		//DbQuery = Db
		Collector.Schema = "postgres."
		Collector.UUID = "('{'||md5(random()::text || clock_timestamp()::text)::uuid||'}')"
		//DbTimeStamp = "(CAST (to_char(now(), 'J') AS INT) - 2440587.5)*86400.0*1000"
		Collector.DbTimeStamp = "(now())"
		log.Print("Pinging Postgresql: ")
		log.Println(Collector.DatabaseDB.Ping)
	}

	/*
		if DbSource == 0 {
			tmpSrc := os.Getenv("DB_SOURCE")
			if len(tmpSrc) > 0 {
				if tmpSrc == "PGSQL" {
					DbSource = PGSQL
				} else if tmpSrc == "SQLITE" {
					DbSource = SQLITE3
				} else {
					DbSource = FILE
				}
			} else {
				DbSource = FILE
			}
		}
	*/
	/*
		if len(DbName) > 0 {
			if DbSource == PGSQL {
				Db, err = sql.Open("postgres", DbName)
				if err != nil {
					log.Fatal(err)
				}
				DbQuery = Db
				Schema = "postgres."
				UUID = "('{'||md5(random()::text || clock_timestamp()::text)::uuid||'}')"
				//DbTimeStamp = "(CAST (to_char(now(), 'J') AS INT) - 2440587.5)*86400.0*1000"
				DbTimeStamp = "(now())"

				log.Print("Postgresql database: " + DbName)
				log.Print("Pinging Postgresql: ")
				log.Println(Db.Ping)
				LoadConfiguration()
			} else if DbSource == SQLITE3 {
				Schema = ""
				TableSuffix = "_evw"
				UUID = "(select '{'||upper(substr(u,1,8)||'-'||substr(u,9,4)||'-4'||substr(u,13,3)||'-'||v||substr(u,17,3)||'-'||substr(u,21,12))||'}' from ( select lower(hex(randomblob(16))) as u, substr('89ab',abs(random()) % 4 + 1, 1) as v) as foo)"
				DbTimeStamp = "(julianday('now') - 2440587.5)*86400.0*1000"

				//use 2 different sqlite files:
				//1st: contains configuration information and JSON data
				//2nd: contains actual data

				//			initializeStr := `PRAGMA automatic_index = ON;
				//	        PRAGMA cache_size = 32768;
				//	        PRAGMA cache_spill = OFF;
				//	        PRAGMA foreign_keys = ON;
				//	        PRAGMA journal_size_limit = 67110000;
				//	        PRAGMA locking_mode = NORMAL;
				//	        PRAGMA page_size = 4096;
				//	        PRAGMA recursive_triggers = ON;
				//	        PRAGMA secure_delete = ON;
				//	        PRAGMA synchronous = NORMAL;
				//	        PRAGMA temp_store = MEMORY;
				//	        PRAGMA journal_mode = WAL;
				//	        PRAGMA wal_autocheckpoint = 16384;
							`

				//log.Println(initializeStr)
				//initializeStr = "PRAGMA synchronous = OFF;PRAGMA cache_size=100000;PRAGMA journal_mode=WAL;"
				//log.Println(initializeStr)

				//Db, err = sql.Open("sqlite3", "file:"+DbName+"?PRAGMA journal_mode=WAL")

				once.Do(initDB)

				Db, err = sql.Open("sqlite3", DbName+SqlFlags)
				if err != nil {
					log.Fatal(err)
				}
				err = Db.Ping()
				if err != nil {
					log.Fatalf("Error on opening database connection: %s", err.Error())
				}

				//&sqlite3.SQLiteConn.LoadExtension("stgeometry_sqlite", "sqlite3_stgeometrysqlite_init")

				//sqlite3.LoadExtension("stgeometry_sqlite", "sqlite3_stgeometrysqlite_init")


				//   conn := &SQLiteConn{db: Db, loc: loc, txlock: txlock}
				//   conn.LoadExtensions()
				//   	if len(d.Extensions) > 0 {
				//   		if err := conn.loadExtensions(d.Extensions); err != nil {
				//   			conn.Close()
				//   			return nil, err
				//   		}
				//   	}


				//_, err = DbQuery.Exec("SELECT load_extension('stgeometry_sqlite')")
				//_, err = DbQuery.("stgeometry_sqlite","sqlite3_stgeometrysqlite_init")
				//sqlite3conn := []*sqlite3.SQLiteConn{}
				//c *sqlite3.SQLiteConn

				//   sql.Register("sqlite3_with_extensions", &sqlite3.SQLiteDriver{
				//   		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				//   			return conn.CreateModule("github", &githubModule{})
				//   		},
				//   	})


				//_, err = Db.Exec("SELECT load_extension('stgeometry_sqlite','SDE_SQL_funcs_init')")
				//SELECT load_extension('stgeometry_sqlite.dll','SDE_SQL_funcs_init');
				//if err != nil {
				//	log.Fatalf("Error on loading extension stgeometry_sqlite: %s", err.Error())
				//}

				//Db.Exec(initializeStr)
				log.Println("Sqlite config database: " + DbName)
				//defer Db.Close()
				//Db.SetMaxOpenConns(1)

				LoadConfiguration()
				//get ServiceName
				DbQueryName := DataPath + string(os.PathSeparator) + ServiceName + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + ServiceName + ".geodatabase"

				//DbQuery, err = sql.Open("sqlite3", "file:"+DbQueryName+"?PRAGMA journal_mode=WAL")
				DbQuery, err = sql.Open("sqlite3", DbQueryName+SqlFlags)
				if err != nil {
					log.Fatal(err)
				}
				err = DbQuery.Ping()
				if err != nil {
					log.Fatalf("Error on opening database connection: %s", err.Error())
				}
				log.Println("Sqlite replica database: " + DbQueryName)

				//testQuery()
				//os.Exit(0)



				//defer DbQuery.Close()
				//DbQuery.SetMaxOpenConns(1)
				//log.Print("Sqlite database: " + DbQueryName)
				//DbQuery.Exec(initializeStr)
				//defer db.Close()
			}
		} else {
			if DbSource == PGSQL {
				log.Println("Missing Postgresql connection string:  defaulting to FileSystem")
			} else if DbSource == SQLITE3 {
				log.Println("Missing Sqlite database name:  defaulting to FileSystem")
			}
			DbSource = FILE
		}
	*/
	/*
		else if DbSource == FILE {
			LoadConfigurationFromFile()
		}
	*/

	/*
		Db, err = sql.Open("postgres", "user=postgres DbSource=gis host=192.168.99.100")
		if err != nil {
			log.Fatal(err)
		}
	*/

	//DataPath = RootPath        //+ string(os.PathSeparator)        //+ defaultService + string(os.PathSeparator) + "services" + string(os.PathSeparator)

	/*
		ReplicaPath = DataPath     //+ string(os.PathSeparator)     //+ defaultService + string(os.PathSeparator) + "replicas" + string(os.PathSeparator)
		AttachmentsPath = DataPath //+ string(os.PathSeparator) + ServiceName + string(os.PathSeparator) + "attachments" //+ string(os.PathSeparator)
		UploadPath = DataPath + string(os.PathSeparator) + ServiceName + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "attachments"

		log.Println("Service name: " + ServiceName)
		log.Println("Data path: " + DataPath)
		log.Println("Replica path: " + ReplicaPath)
		log.Println("Attachments path: " + AttachmentsPath)
		var DbSourceName string
		switch DbSource {
		case FILE:
			DbSourceName = "Filesystem"
			break
		case PGSQL:
			DbSourceName = "Postgresql"
			break
		case SQLITE3:
			DbSourceName = "Sqlite"
			break
		default:
			DbSourceName = "Unknown"
		}
		log.Println("Data source: " + DbSourceName)
		log.Println("Data name" + DbName)
	*/

	//print out summary
	PrintServerSummary()
}
func PrintServerSummary() {
	log.Printf("HTTP Port: %v\n", Collector.HttpPort)
	log.Printf("HTTPS Port: %v\n", Collector.HttpsPort)
	log.Printf("Cert: %v\n", Collector.Pem)
	log.Printf("Pem: %v\n", Collector.Cert)
	log.Printf("Sqlite configuration DB %v\n", Collector.SqliteDb)
	for key, _ := range Collector.Projects {
		log.Printf("Loading project:  %v\n", key)
		//log.Printf("%v %v\n", key, val.DataPath)
		//log.Printf("%v %v\n", key, val.FGDB)
		//log.Printf("%v %v\n", key, val.ReplicaPath)

	}
	if Collector.DefaultDataSource == structs.PGSQL {
		log.Println("Using Postgresql database: " + Collector.PG)

	} else if Collector.DefaultDataSource == structs.SQLITE3 {
		log.Println("Using SQLITE replica database ")

	} else if Collector.DefaultDataSource == structs.FILE {
		log.Println("Using File based database")
	}

}

func PrintServerSummaryTable(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Welcome"))
	w.Header().Set("Content-Type", "application/json")
	/*
		var info = map[string]interface{}{} //{"server": server, "projects": projects}

		info["HTTP Port"] = Collector.HttpPort
		info["HTTPS Port"] = Collector.HttpsPort
		info["Cert"] = Collector.Pem
		info["Pem"] = Collector.Cert
		info["Sqlite configuration DB"] = Collector.SqliteDb
		if Collector.DefaultDataSource == structs.PGSQL {
			info["Postgresql database"] = Collector.PG
		} else if Collector.DefaultDataSource == structs.SQLITE3 {

		} else if Collector.DefaultDataSource == structs.FILE {
		}
		var projects = map[string]map[string]interface{}{}
		//info["projects"] = projects
		//var projects map[string]interface{} //{"itemID": "1", "itemName": fileName, "description": "description", "date": time.Now().Local().Unix() * 1000, "committed": true}
		for _, val := range Collector.Projects {
			projects[val.Name] = map[string]interface{}{}
			projects[val.Name]["name"] = val.Name
			projects[val.Name]["datapath"] = val.DataPath
			projects[val.Name]["fgdb"] = val.FGDB
			projects[val.Name]["replicaPath"] = val.ReplicaPath
		}
		info["projects"] = projects
	*/
	//item, _ := json.Marshal(map[string]interface{}{"itemID": "1", "itemName": fileName, "description": "description", "date": time.Now().Local().Unix() * 1000, "committed": true})
	response, _ := json.Marshal(Collector)
	w.Write(response)
}

func initDB() {
	sql.Register("sqlite3_with_extensions",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"stgeometry_sqlite",
			},
		})
}

func GetParam(DbSource string, i int) string {
	if DbSource == structs.SQLITE3 {
		return "?"
	}
	return "$" + strconv.Itoa(i)
}

/*
func SetDatasource(newDatasource int) {
	if DbSource == newDatasource {
		return
	}
	if newDatasource == FILE {
		DbSource = FILE
		//close db
		Db.Close()
		return
	}

	var err error
	if newDatasource == PGSQL {

		Db, err = sql.Open("postgres", Project.PG)
		if err != nil {
			log.Fatal(err)
		}
		DbQuery = Db
		log.Print("Postgresql database: " + Project.PG)
		log.Print("Pinging Postgresql: ")
		log.Println(Db.Ping)
	} else if newDatasource == SQLITE3 {
		Db, err = sql.Open("sqlite3", "file:"+Collector.SqliteDb+SqlWalFlags)
		if err != nil {
			log.Fatal(err)
		}
		DbQueryName := DataPath + string(os.PathSeparator) + ServiceName + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + ServiceName + ".geodatabase"
		log.Println("DbQueryName: " + DbQueryName)
		DbQuery, err = sql.Open("sqlite3", "file:"+DbQueryName+SqlWalFlags)

		if err != nil {
			log.Fatal(err)
		}
	}
}
*/
/*
func LoadConfiguration() {

	sql := "select json from catalog where name=" + GetParam(1)
	log.Printf("Query: select json from catalog where name='config'")
	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	var str []byte
	err = stmt.QueryRow("config").Scan(&str)
	if err != nil {
		log.Println("Error reading configuration table")
		log.Println(err.Error())
		LoadConfigurationFromFile()
		return
	}
	err = json.Unmarshal(str, &Collector)
	if err != nil {
		log.Println("Error parsing configuration table")
		log.Println(err.Error())
		LoadConfigurationFromFile()
		return
	}
	//Project = Collector.Projects[0]
	//for i, _ := range Project.Services {
	//	ServiceName = i
	//	//ServiceName = val[i]
	//}
}
*/
func LoadConfigurationFromFile(DataPath string) {
	var configFile = DataPath + string(os.PathSeparator) + "config.json"
	//var json []byte
	file, err1 := ioutil.ReadFile(configFile)
	if err1 != nil {
		fmt.Printf("// error while reading file %s\n", configFile)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	err2 := json.Unmarshal(file, &Collector)
	if err2 != nil {
		log.Println("Error reading configuration file: " + configFile)
		log.Println(err2.Error())
	}
	//Project = Collector.Projects[0]
	//for i, _ := range Project.Services {
	//	ServiceName = i
	//RootName = val[i]
	//}

}

func GetDataBase(project structs.Project) *sql.DB {
	/*
		if project == nil {
			//Collector.DefaultDataSource
		} else if project.DB == nil {

		}
		//sql := "ATTACH DATABASE '" + deltaFile + "' AS delta"
	*/
	return nil
}

func GetReplicaDB(name string) *sql.DB {
	if Collector.Projects[name].ReplicaDB == nil {
		var err error
		//var db *sql.DB
		//Collector.Projects[name].ReplicaDB = make(sql.DB)
		//var p structs.Project = Collector.Projects[name]
		//Collector.Projects[name].ReplicaDB = new(sql.DB)
		Collector.Projects[name].ReplicaDB, err = sql.Open("sqlite3", Collector.Projects[name].ReplicaPath+SqlFlags)
		if err != nil {
			log.Fatal(err)
		}
		err = Collector.Projects[name].ReplicaDB.Ping()
		if err != nil {
			log.Fatalf("Error on opening database connection: %s", err.Error())
		}
		//Collector.Projects[name].ReplicaDB = p.ReplicaDB
	}
	return Collector.Projects[name].ReplicaDB
}

/*
func GetReplicaDB(name string) *sql.DB {
	if Collector.ReplicaDB[name] == nil {
		var err error
		//var db *sql.DB
		//Collector.Projects[name].ReplicaDB = make(sql.DB)
		Collector.ReplicaDB[name], err = sql.Open("sqlite3", Collector.Projects[name].ReplicaPath+SqlFlags)
		if err != nil {
			log.Fatal(err)
		}
		err = Collector.ReplicaDB[name].Ping()
		if err != nil {
			log.Fatalf("Error on opening database connection: %s", err.Error())
		}

	}
	return Collector.ReplicaDB[name]
}
*/
//GetArcService queries the database for service layer entries
func GetArcService(catalog string, service string, layerid int, dtype string, dbPath string) []byte {

	if Collector.DefaultDataSource == structs.FILE {
		if len(service) > 0 {
			service += "."
		}
		sp := ""
		if layerid > -1 {
			sp = fmt.Sprint(layerid, ".")
		}

		if len(dtype) > 0 {
			if dtype == "data" && service == "content" {
				dtype = "items." + dtype + "."
			} else {
				dtype += "."
			}

		}
		jsonFile := fmt.Sprint(Collector.DataPath, string(os.PathSeparator), catalog, string(os.PathSeparator), "services", string(os.PathSeparator), service, sp, dtype, "json")
		file, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			log.Println(err)
		}
		return file
	}
	sql := "select json from services where service like " + GetParam(Collector.DefaultDataSource, 1) + " and name=" + GetParam(Collector.DefaultDataSource, 2) + " and layerid=" + GetParam(Collector.DefaultDataSource, 3) + " and type=" + GetParam(Collector.DefaultDataSource, 4)
	log.Printf("Query: select json from services where service like '%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
	//Db := GetDataBase(Collector.Projects[catalog])
	stmt, err := Collector.Configuration.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	var json []byte
	err = stmt.QueryRow(catalog, service, layerid, dtype).Scan(&json)
	if err != nil {
		log.Println(err.Error())
		//log.Println(sql)
	}
	return json
}

//GetArcCatalog queries the database for top level catalog entries
func GetArcCatalog(service string, dtype string, dbPath string) []byte {

	if Collector.DefaultDataSource == structs.FILE || service == "config" {
		if len(service) > 0 {
			service += "."
		}

		if len(dtype) > 0 {
			dtype += "."
		}

		jsonFile := fmt.Sprint(Collector.DataPath, string(os.PathSeparator), service, dtype, "json")
		file, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			log.Println(err)
		}

		return file

	}
	sql := "select json from catalog where name=" + GetParam(Collector.DefaultDataSource, 1) + " and type=" + GetParam(Collector.DefaultDataSource, 2)
	log.Printf("Query: select json from catalog where name='%v' and type='%v'", service, dtype)
	//Db := GetDataBase(Collector.Projects[catalog])

	stmt, err := Collector.Configuration.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}

	var json []byte
	err = stmt.QueryRow(service, dtype).Scan(&json)
	if err != nil {
		log.Println(err.Error())
		//log.Println(sql)
	}

	return json
}

func SetArcService(json []byte, catalog string, service string, layerid int, dtype string, dbPath string) bool {
	if service == "info" {
		return true
	}
	if Collector.DefaultDataSource == structs.FILE {
		if len(service) > 0 {
			service += "."
		}
		sp := ""
		if layerid > -1 {
			sp = fmt.Sprint(layerid, ".")
		}

		if len(dtype) > 0 {
			if dtype == "data" && service == "content" {
				dtype = "items." + dtype + "."
			} else {
				dtype += "."
			}
		}

		jsonFile := fmt.Sprint(Collector.DataPath, string(os.PathSeparator), catalog, string(os.PathSeparator), "services", string(os.PathSeparator), service, sp, dtype, "json")
		err := ioutil.WriteFile(jsonFile, json, 0644)
		if err != nil {
			return false
		}
		return true
	}
	//Db := GetDataBase(Collector.Projects[catalog])
	sql := "update services set json=" + GetParam(Collector.DefaultDataSource, 1) + " where service like " + GetParam(Collector.DefaultDataSource, 2) + " and name=" + GetParam(Collector.DefaultDataSource, 3) + " and layerid=" + GetParam(Collector.DefaultDataSource, 4) + " and type=" + GetParam(Collector.DefaultDataSource, 5)
	log.Printf("Query: update services set json=<json> where service like '%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
	stmt, err := Collector.Configuration.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	//err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
	_, err = stmt.Exec(json, catalog, service, layerid, dtype)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

//GetArcCatalog queries the database for top level catalog entries
func SetArcCatalog(json []byte, service string, dtype string, dbPath string) bool {
	if service == "info" {
		return true
	}

	if Collector.DefaultDataSource == structs.FILE || service == "config" {
		if len(service) > 0 {
			service += "."
		}
		if len(dtype) > 0 {
			dtype += "."
		}

		jsonFile := fmt.Sprint(Collector.DataPath, string(os.PathSeparator), service, dtype, "json")
		err := ioutil.WriteFile(jsonFile, json, 0644)
		if err != nil {
			return false
		}
		return true
	}
	//Db := GetDataBase(Collector.Projects[catalog])
	sql := "update catalog set json=" + GetParam(structs.SQLITE3, 1) + " where name=" + GetParam(structs.SQLITE3, 2) + " and type=" + GetParam(structs.SQLITE3, 3)
	log.Printf("Query: update catalog set json=<json> where name='%v' and type='%v'", service, dtype)

	stmt, err := Collector.Configuration.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}

	_, err = stmt.Exec(json, service, dtype)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func GetArcQuery(catalog string, service string, layerid int, dtype string, oidname string, objectIds string, where string) []byte {
	//objectIdsInt, _ := strconv.Atoi(objectIds)
	objectIdsArr := strings.Split(objectIds, ",")
	var objectIdsFloat = []float64{}
	for _, i := range objectIdsArr {
		j, err := strconv.ParseFloat(i, 64)
		if err != nil {
			panic(err)
		}
		objectIdsFloat = append(objectIdsFloat, j)
	}

	if Collector.DefaultDataSource == structs.FILE {
		//config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json"

		jsonFile := fmt.Sprint(Collector.DataPath, string(os.PathSeparator), catalog+string(os.PathSeparator), "services", string(os.PathSeparator), "FeatureServer.", layerid, ".query.json")
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
		var results []structs.Feature
		for _, i := range srcObj.Features {
			//if int(i.Attributes["OBJECTID"].(float64)) == objectIdsInt {
			oid := i.Attributes[oidname].(float64)
			if in_float_array(oid, objectIdsFloat) {
				//oJoinVal = i.Attributes[oJoinKey]
				results = append(results, i)
				//break
			}
		}
		srcObj.Features = results
		jsonstr, err := json.Marshal(srcObj)
		if err != nil {
			log.Println(err)
		}
		return jsonstr
	} else if Collector.DefaultDataSource == structs.PGSQL {
		sql := "select json from " + Collector.Schema + "services where service=$1 and name=$2 and layerid=$3 and type=$4"
		log.Printf("select json from "+Collector.Schema+"services where service='%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
		//Db := GetDataBase(Collector.Projects[catalog])
		stmt, err := Collector.Configuration.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		var fields []byte
		err = stmt.QueryRow(catalog, service, layerid, dtype).Scan(&fields)
		if err != nil {
			log.Println(err.Error())
			//w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return []byte("")
		}
		var featureObj structs.FeatureTable
		//var fieldsArr []structs.Field
		err = json.Unmarshal(fields, &featureObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		var results []structs.Feature
		for _, i := range featureObj.Features {
			//if int(i.Attributes["OBJECTID"].(float64)) == objectIdsInt {
			oid := i.Attributes[oidname].(float64)
			if in_float_array(oid, objectIdsFloat) {
				//oJoinVal = i.Attributes[oJoinKey]
				results = append(results, i)
				//break
			}
		}
		featureObj.Features = results
		fields, err = json.Marshal(featureObj)
		if err != nil {
			log.Println(err)
		}
		return fields
	} else if Collector.DefaultDataSource == structs.SQLITE3 {
		sql := "select json from services where service=? and name=? and layerid=? and type=?"
		log.Printf("select json from services where service='%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
		//Db := GetDataBase(Collector.Projects[catalog])
		stmt, err := Collector.Configuration.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		var fields []byte
		err = stmt.QueryRow(catalog, service, layerid, dtype).Scan(&fields)
		if err != nil {
			log.Println(err.Error())
			//w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return []byte("")
		}
		var featureObj structs.FeatureTable
		err = json.Unmarshal(fields, &featureObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		var results []structs.Feature
		for _, i := range featureObj.Features {
			//if int(i.Attributes["OBJECTID"].(float64)) == objectIdsInt {

			oid := i.Attributes[oidname].(float64)
			if in_float_array(oid, objectIdsFloat) {
				//oJoinVal = i.Attributes[oJoinKey]
				results = append(results, i)
				//break
			}
		}
		//globalIdFieldName=GlobalID
		//objectIdFieldName=OBJECTID
		//featureObj.GlobalIDField = "GlobalID"
		//featureObj.ObjectIDFieldName = "OBJECTID"
		/*
			GlobalIDField     string `json:"globalIdField,omitempty"`
			GlobalIDFieldName string `json:"globalIdFieldName,omitempty"`
			ObjectIDField      string `json:"objectIdField,omitempty"`
			ObjectIDFieldName string    `json:"objectIdFieldName,omitempty"`
		*/

		featureObj.Features = results
		fields, err = json.Marshal(featureObj)
		if err != nil {
			log.Println(err)
		}
		return fields
	}
	return []byte("")
	/*
		sql := "select * from  set json=" + GetParam(1) + " where name=" + GetParam(2) + " and type=" + GetParam(3)
		log.Printf("Query: update catalog set json=<json> where name='%v' and type='%v'", service, dtype)

		stmt, err := Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		_, err = stmt.Exec(json, service, dtype)
		if err != nil {
			log.Println(err.Error())
			return false
		}

		return true
	*/
}

func in_string_array(val string, array []string) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}

func in_float_array(val float64, array []float64) bool {
	for i := range array {
		if array[i] == val {
			return true
		}
	}
	return false
}

func in_array(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}

func DblQuote(s string) string {
	return "\"" + s + "\""
}

/*
func testQuery() {

	sql := "insert into grazing_inspections(yearling_heifers,studs,lambs,wethers,kids,reviewer_name,reviewer_date,reviewer_title,cows,steer_calves,mares,fillies,nannies,Comments,OBJECTID,colts,ewes,rams,billies,yearling_steers,bulls,geldings,GlobalGUID,GlobalID) values( ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	vals := []interface{}{nil, nil, nil, nil, nil, nil, nil, nil, 13, nil, nil, nil, nil, nil, 20, nil, nil, nil, nil, nil, nil, nil, "{6FC17403-5889-4A23-AC77-3B060E4C6DC4}", "{6FC17403-5889-4A23-AC77-3B060E4C6DC4}"}
	stmt, err := DbQuery.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	_, err = stmt.Exec(vals...)
	if err != nil {
		log.Println(err.Error())
	}
	stmt.Close()

}
*/
/*
	sql := "select * from grazing_inspections where GlobalGUID in (select substr(GlobalID, 2, length(GlobalID)-2) from grazing_permittees where OBJECTID in(?))"
	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query(16) //relationshipIdInt
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//for i, _ := range columns {
	//	log.Println(columns[i])
	//}

	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
			log.Println(i)
		}
		rows.Scan(valuePtrs...)
		for i, col := range columns {
			log.Println(i)
			log.Println(col)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()

	sql = "select * from grazing_inspections where GlobalGUID in (select substr(GlobalID, 2, length(GlobalID)-2) from grazing_permittees where OBJECTID in(16))"
	sql = "select substr(GlobalID, 2, length(GlobalID)-2) as GlobalGUID from grazing_permittees where OBJECTID in(16)"
	sql = "select OBJECTID,cows,yearling_heifers,steer_calves,yearling_steers,bulls,mares,geldings,studs,fillies,colts,ewes,lambs,rams,wethers,kids,billies,nannies,Comments,GlobalGUID,created_user,created_date,last_edited_user,last_edited_date,reviewer_name,reviewer_date,reviewer_title,GlobalID from grazing_inspections"
	sql = "select * from grazing_permittees"
	sql = "select OBJECTID from grazing_inspections"

	log.Println(sql)
	rows, err = Db.Query(sql) //relationshipIdInt
	columns, _ = rows.Columns()
	count = len(columns)
	values = make([]interface{}, count)
	valuePtrs = make([]interface{}, count)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
			log.Println(i)
		}
		rows.Scan(valuePtrs...)
		for i, col := range columns {
			log.Println(i)
			log.Println(col)
			//val := values[i]
			//log.Printf("%v", val.([]uint8))
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()

	os.Exit(1)
*/
