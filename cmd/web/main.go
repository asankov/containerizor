package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type application struct {
	log *log.Logger
}

func main() {
	port := flag.Int("port", 4000, "port on which the application is exposed")
	flag.Parse()

	app := &application{
		log: log.New(os.Stdout, "", log.Ldate),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/containers", app.listContainers)
	mux.HandleFunc("/containers/create", app.createContainer)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}
	log.Printf("Listening on port %d", *port)
	log.Fatal(srv.ListenAndServe())
}
