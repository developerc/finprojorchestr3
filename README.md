создал папку finprojorchestr3 для финального проекта
установил плагины
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
C:\proj\Go\Exam\finprojorchestr2> go mod init github.com/developerc/finprojorchestr3
создал папку proto и grpc
mkdir proto
mkdir grpc
в папке proto создал файл grpc.proto
создал файлы 
grpc/client/clnt.go
grpc/server/srv.go
залил на github
сгенерировал файлы
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     proto/grpc.proto  
go mod tidy
cd grpc
C:\proj\Go\Exam\finprojorchestr3\grpc> go mod init grpc
заполнил srv.go и подтянул зависимости
go mod tidy
replace grpc => ./grpc/
go get grpc
в главной ветке создадим main.go
mkdir http
go mod init http
заполнили srv.go подтянули зависимости
go mod tidy
replace grpc => ../grpc/
go get grpc
go mod tidy
создал папку sqlite
mkdir sqlite
cd sqlite
go mod init sqlite
go mod tidy
заполнили srv.go подтянули зависимости
go get github.com/mattn/go-sqlite3
go mod tidy
в файле go.mod
replace sqlite => ./sqlite/
go get sqlite
go mod tidy
//---
создал db.go
cd grpc
go get github.com/mattn/go-sqlite3
go mod tidy

sqlite3 store.db
sqlite> .tables
sqlite> select * from tasks;
sqlite> .quit
//----