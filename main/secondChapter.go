package main

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	s "WebApplication/main/secondary"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var DB_NAME = "mysql"
var DB_CREDENTIALS = "root:i86Kwp5a@tcp(127.0.0.1:3306)/new_schema"
var ENDPOINT = "/api/user"
var LOCAL_HOST = ":8080"

func main() {
	router := mux.NewRouter()
	router.HandleFunc(ENDPOINT, CreateUser).Methods("POST")
	router.HandleFunc(ENDPOINT, RetrieveUsers).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(LOCAL_HOST, nil)
}

/**
reading query parameters, creating User object and writing it to mySql database
*/
func CreateUser(responseWriter http.ResponseWriter, request *http.Request) {
	User := s.User{}
	User.ID, _ = strconv.Atoi(request.FormValue("id"))
	User.Name = request.FormValue("name")
	User.Email = request.FormValue("email")
	User.First = request.FormValue("first")
	User.Last = request.FormValue("last")
	_, error := json.Marshal(User)
	if error != nil {
		fmt.Println("marshalling error!")
	}

	database, error := sql.Open(DB_NAME, DB_CREDENTIALS)
	defer database.Close()

	if error != nil {
		fmt.Println("couldn't open db")
	}

	_, err := database.Exec(
		"INSERT INTO users set ID='" + strconv.Itoa(User.ID) +
			"', Nickname='" + User.Name +
			"', First='" + User.First +
			"', Last='" + User.Last +
			"', Email='" + User.Email +
			"'")
	if err != nil {
		fmt.Fprintf(responseWriter, err.Error())
	}
	fmt.Fprintf(responseWriter, "object written to database")
	fmt.Fprintf(responseWriter, User.ToString())
}

/**
Read a user from sql database via webapplication api
*/
func ReadUser(responseWriter http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	ReadUser := s.User{}

	database, error := sql.Open(DB_NAME, DB_CREDENTIALS)
	defer database.Close()
	if error != nil {
		fmt.Println("couldn't open db")
	}

	err := database.QueryRow("select * from users where ID=?", id).Scan(&ReadUser.ID, &ReadUser.Name, &ReadUser.First, &ReadUser.Last, &ReadUser.Email)
	if err != nil {
		fmt.Fprintf(responseWriter, err.Error())
	}
	answer, _ := json.Marshal(ReadUser)
	fmt.Fprintf(responseWriter, string(answer))
}

func RetrieveUsers(responseWriter http.ResponseWriter, request *http.Request) {
	response := s.Users{}
	responseWriter.Header().Set("Pragma", "no-cache")
	database, openErr := sql.Open(DB_NAME, DB_CREDENTIALS)
	defer database.Close()
	if openErr != nil {
		fmt.Fprintf(responseWriter, openErr.Error())
	}
	rows, readErr := database.Query("select * from users LIMIT 10")
	if readErr != nil {
		fmt.Fprintf(responseWriter, readErr.Error())
	}
	for rows.Next() {
		user := s.User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
		response.Users = append(response.Users, user)
	}
	jsonResponse, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		fmt.Fprintf(responseWriter, marshalErr.Error())
	}
	fmt.Fprintf(responseWriter, string(jsonResponse))
}
