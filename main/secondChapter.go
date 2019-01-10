package main

import (
	s "WebApplication/main/secondary"
	"database/sql"
	_ "database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

const DB_NAME = "mysql"
const DB_CREDENTIALS = "root:i86Kwp5a@tcp(127.0.0.1:3306)/new_schema"
const ENDPOINT = "/api/user"
const LOCAL_HOST = ":8080"

var Database *sql.DB

func main() {
	db, error := sql.Open(DB_NAME, DB_CREDENTIALS)
	if error != nil {
		fmt.Println("couldn't open db")
	}
	defer db.Close()
	Database = db

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
	file, _, _ := request.FormFile("image")
	fileData, _ := ioutil.ReadAll(file)
	fileString := base64.StdEncoding.EncodeToString(fileData)
	User.Image = fileString
	_, error := json.Marshal(User)
	if error != nil {
		fmt.Println("marshalling error!")
	}

	_, err := Database.Exec(
		"INSERT INTO users set ID='" + strconv.Itoa(User.ID) +
			"', Nickname='" + User.Name +
			"', First='" + User.First +
			"', Last='" + User.Last +
			"', Email='" + User.Email +
			"', Image='" + User.Image +
			"'")
	if err != nil {
		fmt.Fprintf(responseWriter, err.Error())
	}
	fmt.Fprintf(responseWriter, "object written to database")
	fmt.Fprintf(responseWriter, User.ToString())
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
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email, &user.Image)
		response.Users = append(response.Users, user)
	}
	jsonResponse, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		fmt.Fprintf(responseWriter, marshalErr.Error())
	}
	fmt.Fprintf(responseWriter, string(jsonResponse))
}

/**
Read a user from sql database via webapplication api
*/
func ReadUser(responseWriter http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	ReadUser := s.User{}

	err := Database.QueryRow("select * from users where ID=?", id).Scan(&ReadUser.ID, &ReadUser.Name, &ReadUser.First, &ReadUser.Last, &ReadUser.Email)
	if err != nil {
		fmt.Fprintf(responseWriter, err.Error())
	}
	answer, _ := json.Marshal(ReadUser)
	fmt.Fprintf(responseWriter, string(answer))
}
