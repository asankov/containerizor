package main

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/asankov/containerizor/pkg/containers"
	"github.com/asankov/containerizor/pkg/models"
)

type templateArgs struct {
	Containers []*containers.Container
	User       *models.User
}

func (srv *server) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
		}

		srv.serveTemplate(w, nil, "./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func (srv *server) handleContainersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, err := srv.userFromRequest(r)
		if err != nil {
			http.Error(w, "unathorized", http.StatusUnauthorized)
			return
		}

		containers, err := srv.orchestrator.ListContainers()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		t, err := template.ParseFiles("./ui/html/list.page.tmpl", "./ui/html/base.layout.tmpl")
		if err != nil {
			srv.log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		if err := t.Execute(w, templateArgs{
			Containers: containers,
			User:       usr,
		}); err != nil {
			srv.log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (srv *server) handleContainersStartView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, err := srv.userFromRequest(r)
		if err != nil {
			http.Error(w, "unathorized", http.StatusUnauthorized)
			return
		}

		srv.serveTemplate(w, &templateArgs{User: usr}, "./ui/html/start.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func (srv *server) handleContainersStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imageName := r.PostFormValue("image")
		if imageName == "" {
			http.Error(w, "Empty image", 400)
			return
		}

		if _, err := srv.orchestrator.StartNewFrom(imageName); err != nil {
			srv.log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		redirectToView(w, "/containers")
	}
}

func (srv *server) serveTemplate(w http.ResponseWriter, data interface{}, templates ...string) {
	t, err := template.ParseFiles(templates...)
	if err != nil {
		srv.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err := t.Execute(w, data); err != nil {
		srv.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (srv *server) handleContainerStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		if err := srv.orchestrator.StopContainer(id); err != nil {
			srv.log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		redirectToView(w, "/containers")
	}
}

func (srv *server) handleContainerStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		if err := srv.orchestrator.StartContainer(id); err != nil {
			srv.log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		redirectToView(w, "/containers")
	}
}

func (srv *server) handleContainerExecView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, err := srv.userFromRequest(r)
		if err != nil {
			http.Error(w, "unathorized", http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
		id := params["id"]

		srv.serveTemplate(w, execContainerViewResult{ID: id, User: usr}, "./ui/html/exec.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

type execContainerViewResult struct {
	ID     string
	Result *containers.ExecResult
	Cmd    string
	User   *models.User
}

func (srv *server) handleContainerExec() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, err := srv.userFromRequest(r)
		if err != nil {
			http.Error(w, "unathorized", http.StatusUnauthorized)
			return
		}

		command := r.PostFormValue("command")
		if command == "" {
			http.Error(w, "Command cannot be empty", 400)
		}

		params := mux.Vars(r)
		id := params["id"]

		execResult, err := srv.orchestrator.ExecIntoContainer(id, command)
		if err != nil {
			srv.log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		srv.serveTemplate(w, execContainerViewResult{ID: id, Result: execResult, Cmd: command, User: usr}, "./ui/html/exec.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func redirectToView(w http.ResponseWriter, url string) {
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (srv *server) userFromRequest(r *http.Request) (*models.User, error) {
	username := r.Header.Get("user")
	if username == "" {
		return nil, errors.New("no user present in request")
	}

	return srv.db.GetUserByUsername(username)
}
