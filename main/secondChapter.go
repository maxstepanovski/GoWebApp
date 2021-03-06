package main

import (
	s "WebApplication/main/secondary"
	"database/sql"
	_ "database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
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
var Format string

func main() {
	StartServer()
}

func StartServer() {
	db, error := sql.Open(DB_NAME, DB_CREDENTIALS)
	HandleError(error)
	defer db.Close()
	Database = db

	router := mux.NewRouter()
	router.HandleFunc(ENDPOINT, CreateUser).Methods("POST")
	router.HandleFunc(ENDPOINT, RetrieveUsers).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(LOCAL_HOST, nil)
}

func GetFormat(request *http.Request) {
	Format = request.FormValue("format")
}

func SetFormat(data interface{}) []byte {
	var apiOutput []byte
	if Format == "json" {
		output, err := json.Marshal(data)
		HandleError(err)
		apiOutput = output
	} else if Format == "xml" {
		output, err := xml.Marshal(data)
		HandleError(err)
		apiOutput = output
	}
	return apiOutput
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

	_, err := Database.Exec(
		"INSERT INTO users set ID='" + strconv.Itoa(User.ID) +
			"', Nickname='" + User.Name +
			"', First='" + User.First +
			"', Last='" + User.Last +
			"', Email='" + User.Email +
			"', Image='" + User.Image +
			"'")
	HandleError(err)
	fmt.Fprintf(responseWriter, "object written to database")
	fmt.Fprintf(responseWriter, User.ToString())
}

func RetrieveUsers(responseWriter http.ResponseWriter, request *http.Request) {
	GetFormat(request)
	response := s.Users{}
	responseWriter.Header().Set("Pragma", "no-cache")

	rows, readErr := Database.Query("select * from users LIMIT 10")
	HandleError(readErr)
	for rows.Next() {
		user := s.User{}
		scanErr := rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email, &user.Image)
		HandleError(scanErr)
		response.Users = append(response.Users, user)
	}
	finalResult := SetFormat(response)
	fmt.Fprintf(responseWriter, string(finalResult))
}

/**
Read a user from sql database via webapplication api
*/
func ReadUser(responseWriter http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	ReadUser := s.User{}

	err := Database.QueryRow("select * from users where ID=?", id).Scan(&ReadUser.ID, &ReadUser.Name, &ReadUser.First, &ReadUser.Last, &ReadUser.Email)
	HandleError(err)
	answer, _ := json.Marshal(ReadUser)
	fmt.Fprintf(responseWriter, string(answer))
}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
