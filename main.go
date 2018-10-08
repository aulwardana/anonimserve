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

	Run(r, "127.0.0.1", 8080)
}
