package http

import (
	"encoding/json"
	"fmt"
	"grpc/server"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

const hmacSampleSecret = "super_secret_signature"

func makeToken() string {

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": "user_name",
		//"nbf":  now.Add(time.Minute).Unix(),
		"nbf": now.Unix(),
		"exp": now.Add(60 * time.Minute).Unix(),
		"iat": now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		panic(err)
	}

	fmt.Println(tokenString)
	return tokenString
}

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

// Обработчик тестового запроса.
func login(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Write([]byte("Ответ от Orchestrator"))
	//http.Redirect(w, r, "https://dzen.ru", http.StatusSeeOther)
	//http.Redirect(w, r, "/site/", http.StatusSeeOther)
	//Методом POST передается выражение для вычисления
	if r.Method != http.MethodPost { //если это не POST
		log.Println("method is no POST!")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("StatusBadRequest"))
		return
	}
	lgn := r.URL.Query().Get("lgn")
	psw := r.URL.Query().Get("psw")
	fmt.Println("lgn", lgn)
	fmt.Println("psw", psw)
	var auth Auth
	auth.Result = "success"
	token := makeToken()
	auth.Token = token
	//В ответе отсылаем результат аутентификации и токен
	js, err := json.Marshal(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func RunHttpSrv() {
	/*fmt.Println("running http server ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)                 //проверка соединения с сервером
	mux.HandleFunc("/send_expr/", handleExpr) //POST Запрос отправки вычисления выражения
	http.ListenAndServe(":8080", mux)*/
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/api/v1/login/", login)
	http.HandleFunc("/api/v1/send_expr/", handleExpr) //POST Запрос отправки вычисления выражения

	log.Print("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
