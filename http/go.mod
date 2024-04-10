module http

go 1.22.1

replace grpc => ../grpc/

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	grpc v0.0.0-00010101000000-000000000000
)

require (
	github.com/developerc/finprojorchestr3 v0.0.0-20240402193307-4505ac6af433 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/grpc v1.62.1 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
