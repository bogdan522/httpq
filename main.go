package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	fmt.Println("Hello World")
	//h := &HTTPQ{}
	setupRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
