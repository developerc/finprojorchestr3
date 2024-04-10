package http

import (
	"encoding/json"
	"fmt"
	"grpc/server"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

type Registration struct {
	Result string `json:"result"`
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

// обработчик запроса на регистрацию
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //если это не POST
		log.Println("method is no POST!")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("StatusBadRequest"))
		return
	}

	lgn := r.URL.Query().Get("lgn")
	psw := r.URL.Query().Get("psw")
	fmt.Println("register lgn", lgn)
	fmt.Println("register psw", psw)
	var reg Registration
	//добавляем в БД пользователя
	err := server.InsertUser(lgn, psw)
	if err != nil {
		//fmt.Println(err)
		if strings.Contains(string(err.Error()), "UNIQUE constraint failed") {
			reg.Result = "UNIQUE constraint failed"
		} else {
			reg.Result = "not success"
		}
	} else {
		reg.Result = "success"
	}
	fmt.Println("err InsertUser: ", err)
	//token := makeToken()
	//auth.Token = token
	//В ответе отсылаем результат регистрации
	js, err := json.Marshal(reg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

// Обработчик запроса на авторизацию
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
	//если аутентификация успешна
	if err := server.IsPswValid(lgn, psw); err == nil {
		auth.Result = "success"
		token := makeToken()
		auth.Token = token
	} else {
		fmt.Println("err: ", string(err.Error()))
		auth.Result = string(err.Error())
		auth.Token = ""
	}

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
	http.HandleFunc("/api/v1/register/", register)
	http.HandleFunc("/api/v1/send_expr/", handleExpr) //POST Запрос отправки вычисления выражения

	log.Print("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
