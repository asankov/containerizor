package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *server) routes() *mux.Router {
	router := mux.NewRouter()
	// POST should be allowed here, because of the redirects
	router.HandleFunc("/", app.home()).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/containers", app.handleContainersList()).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/containers/start", app.handleContainersStartView()).Methods(http.MethodGet)
	router.HandleFunc("/containers/start", app.handleContainersStart()).Methods(http.MethodPost)
	router.HandleFunc("/containers/{id}/stop", app.handleContainerStop()).Methods(http.MethodPost)
	router.HandleFunc("/containers/{id}/start", app.handleContainerStart()).Methods(http.MethodPost)
	router.HandleFunc("/containers/{id}/exec", app.handleContainerExecView()).Methods(http.MethodGet)
	router.HandleFunc("/containers/{id}/exec", app.handleContainerExec()).Methods(http.MethodPost)

	router.HandleFunc("/users/create", app.handleUserCreateView()).Methods(http.MethodGet)
	router.HandleFunc("/users/create", app.handleUserCreate()).Methods(http.MethodPost)
	router.HandleFunc("/users/login", app.handleLoginView()).Methods(http.MethodGet)
	router.HandleFunc("/users/login", app.handleLogin()).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return router
}
