# collector-server

Collector server is a replacement for using ArcGIS Server Online to host, load, and edit simple data layers.  Written using Go, it has options for using the filesystem, Sqlite, or Postgresql (with Postgis) as the backend database.

## Getting Started
This is only intended to work on simple data layers for use with the ESRI Collector App.  The Collector App connects to the server and shows the map containing 

Use https://github.com/traderboy/collector-tools to create the database to use for the server.

### Prerequisites
* Windows or Linux OS
* Golang 1.7.1 (or later, may work with earlier versions too)
* (Optional) Postgresql 10+ with Postgis

### Installing
Requires the following go libraries installed:
````
go get github.com/gorilla/handlers
go get github.com/lib/pq
go get github.com/gorilla/mux
go get github.com/mattn/go-sqlite3
go get github.com/twinj/uuid
````

### Running
First, create your database using collector-tools.  The data should be created in a "catalogs" directory

To view the various configuration files, open your browser and go to http://localhost.  You can customize the display fields as well as set layers to editable/non-editable.

To add/edit shapes, you must have sqlite3.exe and the following .dlls from ArcMap installed in the root folder of the server executable:
stgeometry_sqlite.dll
icudt52.dll
icuio52.dll
icuin52.dll
icuuc52.dll


### Docker
Docker instructions
````
docker build -t traderboy/collector-server -f docker/Dockerfile .
````

````
docker run -d -p 80:80 -p 443:443 -e HTTPS_PORT=443 -e HTTP_PORT=80 -e ROOT_PATH=catalogs -e DB_SOURCE=PGSQL -e DB_NAME="user=postgres dbname=gis host=172.17.0.5" --name collector-server traderboy/collector-server
````

Helpful to see the docker output while server is running
````
docker logs collector-server -f
````


