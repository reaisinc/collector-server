#FROM postgres:9.5
#kartoza/postgis:9.5-2.2
#FROM golang:onbuild

#   docker stop collector-server;docker rm collector-server;docker build -t traderboy/collector-server -f docker/Dockerfile .;docker run -d -p 80:80 -p 443:443 -e HTTPS_PORT=443 -e HTTP_PORT=80 -e ROOT_PATH=catalogs -e DB_SOURCE=PG -e DB_NAME="user=postgres dbname=gis host=172.17.0.5" --name collector-server traderboy/collector-server;docker logs collector-server -f
#   docker build –t traderboy/collector-server -f docker/Dockerfile .
#   docker rm collector-server
#   docker run -d -p 80:80 -p 443:443 -e HTTPS_PORT=443 -e HTTP_PORT=80 --name collector-server traderboy/collector-server
#   docker run -d -p 80:8000 --name collector-server --link db:postgres traderboy/collector-server
#   docker logs collector-server
#   docker rm collector-server

FROM golang:1.7
#RUN mkdir /app
#RUN mkdir /usr/local/go/src/github.com/traderboy
RUN mkdir -p /usr/local/go/src/github.com/traderboy/collector-server
RUN go get github.com/gorilla/handlers
RUN go get github.com/lib/pq
RUN go get github.com/gorilla/mux
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/twinj/uuid


ADD structs /usr/local/go/src/github.com/traderboy/collector-server/structs
ADD routes /usr/local/go/src/github.com/traderboy/collector-server/routes
ADD config /usr/local/go/src/github.com/traderboy/collector-server/config
#ADD catalogs /usr/local/go/src/github.com/traderboy/collector-server/catalogs
ADD dist /usr/local/go/src/github.com/traderboy/collector-server/dist
ADD ssl /usr/local/go/src/github.com/traderboy/collector-server/ssl
#ADD sha256 /usr/local/go/src/github.com/traderboy/collector-server/sha256
#ADD arcrest.sqlite /usr/local/go/src/github.com/traderboy/collector-server/

ADD server.go /usr/local/go/src/github.com/traderboy/collector-server 
#ADD docker/server.go /usr/local/go/src/github.com/traderboy/collector-server 
#ADD docker/config.go /usr/local/go/src/github.com/traderboy/collector-server/config
ADD index.html /usr/local/go/src/github.com/traderboy/collector-server

RUN chmod 775 /usr/local/go/src/github.com/traderboy/collector-server/index.html
#RUN chmod 775 /usr/local/go/src/github.com/traderboy/collector-server/arcrest.sqlite
RUN chmod -R 775 /usr/local/go/src/github.com/traderboy/collector-server/dist


#RUN ls -l /usr/local/go/src/github.com/traderboy/collector-server
#RUN ls -l /usr/local/go/src/github.com/traderboy/collector-server/config

WORKDIR /usr/local/go/src/github.com/traderboy/collector-server
# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
#RUN go install github.com/traderboy/collector-server
RUN go build /usr/local/go/src/github.com/traderboy/collector-server/server.go

#CMD ["./server","-root catalogs","-file"]
#CMD ["./server","-root","catalogs","-sqlite","collector-server.sqlite"]
#CMD ["./server","-root","catalogs","-pgsql","user=postgres dbname=gis host=172.17.0.5"]
CMD ["./server"]

#ADD myapp /app/
#ENTRYPOINT ["./server"]

#EXPOSE 8080
