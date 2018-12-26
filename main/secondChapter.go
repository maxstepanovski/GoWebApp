package main

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	s "firstChapter/main/secondary"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var DB_NAME = "mysql"
var DB_CREDENTIALS = "root:i86Kwp5a@tcp(127.0.0.1:3306)/new_schema"
var CREATE_ENDPOINT = "/api/user/create"
var READ_ENDPOINT = "/api/user/read"
var LOCAL_HOST = ":8080"

func main() {
	router := mux.NewRouter()
	router.HandleFunc(CREATE_ENDPOINT, CreateUser)
	router.HandleFunc(READ_ENDPOINT, ReadUser)
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
	fmt.Fprintf(responseWriter, ReadUser.ToString())
}
