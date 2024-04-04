package http

import (
	"fmt"
	"grpc/server"
	"net/http"
)

// Обработчик тестового запроса.
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("Ответ от Orchestrator"))
	//server.SndTsk()
}

func handleExpr(w http.ResponseWriter, r *http.Request) { //обрабатываем принятый запрос с выражением
	//Методом POST передается выражение для вычисления
	if r.Method != http.MethodPost { //если это не POST
		fmt.Println("method is no POST!")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("StatusBadRequest"))
		return
	}
	expr := r.URL.Query().Get("expr")
	fmt.Println(expr)
	//server.SndTsk(expr)
	server.HandleHttpExpr(expr)
}

func RunHttpSrv() {
	fmt.Println("running http server ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)                 //проверка соединения с сервером
	mux.HandleFunc("/send_expr/", handleExpr) //POST Запрос отправки вычисления выражения
	http.ListenAndServe(":8080", mux)
}
