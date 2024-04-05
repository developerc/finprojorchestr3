package main

import (
	"grpc/server"
	"http"
	"sqlite"
)

func main() {
	sqlite.CreateSqliteDb()
	go http.RunHttpSrv()
	server.CreateOrchGRPCserver()
}
