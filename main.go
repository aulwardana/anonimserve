package main

import (
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	Route(r, "/", MainHome, GET)
	Route(r, "/error/", errorPageHandler, GET)
	Route(r, "/upload/", uploadHandler, GET)
	Route(r, "/anonim/", UploadFile, POST)

	Run(r, "localhost", 8080)
}
