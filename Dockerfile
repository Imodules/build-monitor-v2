FROM golang:1.8.3-alpine

RUN mkdir /go/src/build-monitor-v2
RUN mkdir /go/src/build-monitor-v2/client
RUN mkdir /go/src/build-monitor-v2/server

WORKDIR /go/src/build-monitor-v2

COPY ./client/dist client
COPY ./server server

WORKDIR server

ENV BM_CLIENT_PATH="../client"
RUN go build -o buildMonitorServer

EXPOSE 3030

ENTRYPOINT ["./buildMonitorServer"]
