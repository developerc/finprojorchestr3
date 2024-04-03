package main

import (
	"grpc/server"
	"http"
)

func main() {
	go http.RunHttpSrv()
	server.CreateOrchGRPCserver()
}
