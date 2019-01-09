package main

import (
	"encoding/json"
	s "WebApplication/main/secondary"
	"fmt"
	"github.com/drone/routes"
	"net/http"
)

func main() {
	muxRouter := routes.New()
	muxRouter.Get("/api/:name", HandleFunction)
	http.Handle("/", muxRouter)
	http.ListenAndServe(":8080", nil)
}

func HandleFunction(responseWriter http.ResponseWriter, request *http.Request) {
	urlParams := request.URL.Query()
	name :=  urlParams.Get(":name")
	//surname := urlParams[":surname"]
	text := "Hello, " + name
	message := s.API{text}
	response, error := json.Marshal(message)
	if error != nil {
		fmt.Println("error happened!")
	}
	fmt.Fprintf(responseWriter, string(response))
}
