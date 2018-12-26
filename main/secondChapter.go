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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/user/create", CreateUser)
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func CreateUser(responseWriter http.ResponseWriter, request *http.Request) {
	User := s.User{}
	User.Name = request.FormValue("name")
	User.Email = request.FormValue("email")
	User.First = request.FormValue("first")
	User.Last = request.FormValue("last")
	_, error := json.Marshal(User)
	if error != nil {
		fmt.Println("marshalling error!")
	}

	database, error := sql.Open("mysql", "root:i86Kwp5a@tcp(127.0.0.1:3306)/new_schema")
	defer database.Close()

	if error != nil {
		fmt.Println("couldn't open db")
	}

	result, err := database.Exec("INSERT INTO users set ID='" + strconv.Itoa(User.ID) + "', Nickname='" + User.Name + "', First='" + User.First + "', Last='" + User.Last + "', Email='" + User.Email + "'")
	if err != nil {
		fmt.Fprintf(responseWriter, err.Error())
	}
	fmt.Println(result)
}
