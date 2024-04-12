package http

import (
	"encoding/json"
	"fmt"
	"grpc/server"
	"log"
	"net/http"
	"strconv"
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

func makeToken(lgn string) string {

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": lgn,
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

func getIdResult(w http.ResponseWriter, r *http.Request) { //Получение результата по ID задачи
	if r.Method != http.MethodGet { //если это не GET
		fmt.Println("method is no GET!")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("StatusBadRequest"))
		return
	}
	id := r.URL.Query().Get("id")
	fmt.Println("id: ", id)
	//вызовем функцию получения задачи по id
	n, err := strconv.ParseInt(id, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", n, n)
	}
	task, err := server.GetTaskById(n)
	if err != nil {
		fmt.Println("not fount task with ID")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write([]byte("not fount task with ID"))
		return
	}
	fmt.Println(task)
	// В ответе JSON с ID нужной задачи
	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
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

	task, err := server.HandleHttpExpr(expr)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable) //406
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("StatusBadRequest"))
		return
	}
	//В ответе отсылаем ID задачи
	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
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
		token := makeToken(lgn)
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

// Middleware авторизация
func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization") // получаем из заголовка значение параметра Authorization
		fmt.Println("Authorization: ", auth)
		tokenFromString, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
			}
			return []byte(hmacSampleSecret), nil
		})
		if err != nil {
			log.Fatal(err)
		}
		if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
			fmt.Println("user name: ", claims["name"])
			//проверим есть ли такое имя в базе
			if err = server.LoginExists(claims["name"].(string)); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized) // вернём ошибку авторизации
			}
		} else {
			//panic(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized) // вернём ошибку авторизации
		}

		/*if auth != "Basic userid:password" {  // авторизация НЕправильная?
			http.Error( w, "Unauthorized", http.StatusUnauthorized )  // вернём ошибку авторизации
			return
		  }*/

		next.ServeHTTP(w, r) // обрабатываем запрос дальше
	}
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
	http.HandleFunc("/api/v1/send_expr/", Authorization(handleExpr))      //POST Запрос отправки вычисления выражения
	http.HandleFunc("/api/v1/get_id_result/", Authorization(getIdResult)) //POST Запрос отправки вычисления выражения

	log.Print("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
