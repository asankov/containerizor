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

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: app.routes(),
	}
	log.Printf("Listening on port %d", *port)
	log.Fatal(srv.ListenAndServe())
}
