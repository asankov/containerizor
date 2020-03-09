package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/asankov/containerizor/pkg/containers"
)

type templateArgs struct {
	Containers []*containers.Container
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}

	app.serveTemplate(w, nil, "./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
}

func (app *application) listContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := app.orchestrator.ListContainers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.ParseFiles("./ui/html/list.page.tmpl", "./ui/html/base.layout.tmpl")
	if err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = t.Execute(w, templateArgs{
		Containers: containers,
	})
	if err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) startContainerIndex(w http.ResponseWriter, r *http.Request) {
	app.serveTemplate(w, nil, "./ui/html/start.page.tmpl", "./ui/html/base.layout.tmpl")
}

func (app *application) startNewContainer(w http.ResponseWriter, r *http.Request) {
	imageName := r.PostFormValue("image")
	if imageName == "" {
		http.Error(w, "Empty image", 400)
		return
	}

	_, err := app.orchestrator.StartNewFrom(imageName)
	if err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	redirectToView(w, "/containers")
}

func (app *application) serveTemplate(w http.ResponseWriter, data interface{}, templates ...string) {
	t, err := template.ParseFiles(templates...)
	if err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) stopContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if err := app.orchestrator.StopContainer(id); err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	redirectToView(w, "/containers")
}

func (app *application) startContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if err := app.orchestrator.StartContainer(id); err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	redirectToView(w, "/containers")
}

func (app *application) execContainerView(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	app.serveTemplate(w, struct{ ID string }{ID: id}, "./ui/html/exec.page.tmpl", "./ui/html/base.layout.tmpl")
}

func redirectToView(w http.ResponseWriter, url string) {
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
