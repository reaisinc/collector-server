package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/handlers"
	config "github.com/traderboy/collector-server/config"
	routes "github.com/traderboy/collector-server/routes"
)

var logPath = "logfile.txt"

func main() {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logParam := flag.Bool("log", false, "a bool")

	if *logParam {
		InitLog()
		log.Println("Writing log file to : logfile.txt")
	} else {
		log.SetOutput(os.Stdout)
		log.Println("Writing log file to stdOut")
	}
	config.Initialize()
	config.Server = ConfigRuntime()
	r := routes.StartGorillaMux()

	//test with: curl -H "Origin: http://localhost" -H "Access-Control-Request-Method: PUT" -H "Access-Control-Request-Headers: X-Requested-With" -X OPTIONS --verbose http://my.host.com/
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	//os.Getenv("ORIGIN_ALLOWED")
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	fmt.Println("Public URL: " + config.Collector.Hostname)
	//bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	bind := fmt.Sprintf("%s:%s", config.Collector.Hostname, config.Collector.HttpPort)

	//  Start HTTP
	go func() {
		if len(config.Collector.Cert) > 0 && len(config.Collector.Pem) > 0 {
			err := http.ListenAndServeTLS(config.Collector.HttpsPort, config.Collector.Cert, config.Collector.Pem, handlers.CORS(originsOk, headersOk, methodsOk)(r)) //handlers.CORS()(r))
			if err != nil {
				log.Fatal("Unable to start HTTPS server: ", err)
			} else {
				log.Println("Started HTTPS server on port " + config.Collector.HttpsPort)
			}
		} else {
			log.Println("Unable to start HTTPS server.  Usage is limited to web page viewer only since Collector app must use HTTPS connection.")
		}
	}()
	// Apply the CORS middleware to our top-level router, with the defaults.
	//err1 := http.ListenAndServe(config.Collector.HttpPort, handlers.CORS(originsOk, headersOk, methodsOk)(r)) //handlers.CORS()(r))
	err1 := http.ListenAndServe(bind, handlers.CORS(originsOk, headersOk, methodsOk)(r)) //handlers.CORS()(r))
	if err1 != nil {
		//log.Println("HTTP server not started on port " + config.Collector.HttpPort)
		log.Fatal("Unable to start HTTP server: ", err1)

	} else {
		log.Println("Started HTTP server on port " + config.Collector.HttpPort)
	}
	/*
		go func() {
			// Apply the CORS middleware to our top-level router, with the defaults.
			err1 := http.ListenAndServe(config.HTTPPort, handlers.CORS(originsOk, headersOk, methodsOk)(r)) //handlers.CORS()(r))
			if err1 != nil {
				log.Fatal("HTTP server: ", err1)
			} else {
				log.Println("Started HTTP server")
			}
		}()
		if len(config.Cert) > 0 && len(config.Pem) > 0 {
			err := http.ListenAndServeTLS(config.HTTPSPort, config.Cert, config.Pem, handlers.CORS(originsOk, headersOk, methodsOk)(r)) //handlers.CORS()(r))
			if err != nil {
				log.Fatal("HTTPS server: ", err)
			} else {
				log.Println("Started HTTPS server")
			}
		} else {
			log.Println("Unable to start HTTPS server")
		}
	*/

}

func InitLog() {

	var err error
	var f *os.File
	f, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		//fmt.fprintln("error opening file: %v", err)
		fmt.Printf("%v\n", err)
	}
	defer f.Close()
	log.SetOutput(f)
}

//InitDb intialize databases
/*
func InitDb() {
	var err error
	Db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

}
*/

//ConfigRuntime print out configuration details
func ConfigRuntime() string {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	log.Printf("Running with %d CPUs\n", nuCPU)
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("IPv4: ", ipv4)
		}
	}
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	//server = ip
	fmt.Println("Public IP: " + ip)
	return ip
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
