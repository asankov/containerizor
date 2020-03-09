package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", app.home).Methods(http.MethodGet)
	// POST should be allowed here, because of the redirects
	router.HandleFunc("/containers", app.listContainers).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/containers/start", app.startContainerIndex).Methods(http.MethodGet)
	router.HandleFunc("/containers/start", app.startNewContainer).Methods(http.MethodPost)
	router.HandleFunc("/containers/{id}/stop", app.stopContainer).Methods(http.MethodPost)
	router.HandleFunc("/containers/{id}/start", app.startContainer).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return router
}
