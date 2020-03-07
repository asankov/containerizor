package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/asankov/containerizor/pkg/containers"
)

type templateArgs struct {
	Containers []*containers.Container
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}

	app.serveTemplate(w, "./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
}

func (app *application) listContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := app.orchestrator.ListContainers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// if len(containers) == 0 {
	// 	fmt.Fprintf(w, "no containers")
	// }
	// for _, c := range containers {
	// 	fmt.Fprintf(w, "ID: %s, Image: %s", c.ID, c.Image)
	// }

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
