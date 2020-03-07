package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}

	app.serveTemplate(w, "./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
}

func (app *application) listContainers(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	if idString == "" {
		fmt.Fprintln(w, "Listing containers...")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil || id < 1 {
		http.Error(w, "Not Found", 404)
		return
	}

	fmt.Fprintf(w, "Showing container with ID %d", id)
}

func (app *application) startContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.serveTemplate(w, "./ui/html/start.page.tmpl", "./ui/html/base.layout.tmpl")
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "Method Now Allowed", 405)
		return
	}

	imageName := r.PostFormValue("image")
	if imageName == "" {
		http.Error(w, "Empty image", 400)
		return
	}

	id, err := app.orchestrator.Start(imageName)
	if err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	fmt.Fprintf(w, "Started a new containers with ID %s", id)
}

func (app *application) serveTemplate(w http.ResponseWriter, templates ...string) {
	t, err := template.ParseFiles(templates...)
	if err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		app.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
