#FROM postgres:9.5
#kartoza/postgis:9.5-2.2
#FROM golang:onbuild

#   docker stop collector;docker rm collector;docker build -t traderboy/collector -f docker/Dockerfile .;docker run -d -p 80:80 -p 443:443 -e HTTPS_PORT=443 -e HTTP_PORT=80 -e ROOT_PATH=catalogs -e DB_SOURCE=PG -e DB_NAME="user=postgres dbname=gis host=172.17.0.5" --name collector traderboy/collector;docker logs collector -f
#   docker build –t traderboy/collector -f docker/Dockerfile .
#   docker rm collector
#   docker run -d -p 80:80 -p 443:443 -e HTTPS_PORT=443 -e HTTP_PORT=80 --name collector traderboy/collector
#   docker run -d -p 80:8000 --name collector --link db:postgres traderboy/collector
#   docker logs collector
#   docker rm collector

FROM golang:1.7
#RUN mkdir /app
#RUN mkdir /usr/local/go/src/github.com/traderboy
RUN mkdir -p /usr/local/go/src/github.com/traderboy/collector
RUN go get github.com/gorilla/handlers
RUN go get github.com/lib/pq
RUN go get github.com/gorilla/mux
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/twinj/uuid


ADD structs /usr/local/go/src/github.com/traderboy/collector/structs
ADD routes /usr/local/go/src/github.com/traderboy/collector/routes
ADD config /usr/local/go/src/github.com/traderboy/collector/config
#ADD catalogs /usr/local/go/src/github.com/traderboy/collector/catalogs
ADD dist /usr/local/go/src/github.com/traderboy/collector/dist
ADD ssl /usr/local/go/src/github.com/traderboy/collector/ssl
#ADD sha256 /usr/local/go/src/github.com/traderboy/collector/sha256
#ADD arcrest.sqlite /usr/local/go/src/github.com/traderboy/collector/

ADD server.go /usr/local/go/src/github.com/traderboy/collector 
#ADD docker/server.go /usr/local/go/src/github.com/traderboy/collector 
#ADD docker/config.go /usr/local/go/src/github.com/traderboy/collector/config
ADD index.html /usr/local/go/src/github.com/traderboy/collector

RUN chmod 775 /usr/local/go/src/github.com/traderboy/collector/index.html
#RUN chmod 775 /usr/local/go/src/github.com/traderboy/collector/arcrest.sqlite
RUN chmod -R 775 /usr/local/go/src/github.com/traderboy/collector/dist


#RUN ls -l /usr/local/go/src/github.com/traderboy/collector
#RUN ls -l /usr/local/go/src/github.com/traderboy/collector/config

WORKDIR /usr/local/go/src/github.com/traderboy/collector
# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
#RUN go install github.com/traderboy/collector
RUN go build /usr/local/go/src/github.com/traderboy/collector/server.go

#CMD ["./server","-root catalogs","-file"]
#CMD ["./server","-root","catalogs","-sqlite","collector.sqlite"]
#CMD ["./server","-root","catalogs","-pgsql","user=postgres dbname=gis host=172.17.0.5"]
CMD ["./server"]

#ADD myapp /app/
#ENTRYPOINT ["./server"]

#EXPOSE 8080
