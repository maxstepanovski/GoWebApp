package main

import (
	_ "database/sql"
	s "firstChapter/main/secondary"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("api/user/create", CreateUser).Methods("GET")
}

func CreateUser(responseWriter http.ResponseWriter, request *http.Request) {
	User := s.User{}
	User.Name = request.FormValue("name")
	User.Email = request.FormValue("email")
	User.First = request.FormValue("first")
	User.Last = request.FormValue("last")
	//marshObject, error := json.Marshal(User)
	//if error != nil {
	//	fmt.Println("marshalling error!")
	//}

}