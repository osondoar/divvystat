package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/osondoar/divvystat/controllers"
	"github.com/osondoar/divvystat/controllers/api"
)

func main() {
	var mainController controllers.MainController
	var loadsController api_controllers.LoadsController

	r := mux.NewRouter()
	r.HandleFunc("/", mainController.Index)
	r.HandleFunc("/api/loads", loadsController.Index)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)

}
