package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", app.home).Methods(http.MethodGet)
	router.HandleFunc("/containers", app.listContainers).Methods(http.MethodGet)
	router.HandleFunc("/containers/start", app.startContainerIndex).Methods(http.MethodGet)
	router.HandleFunc("/containers/start", app.startContainer).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return router
}
