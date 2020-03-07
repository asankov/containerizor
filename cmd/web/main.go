package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asankov/containerizor/pkg/containers"
	"github.com/docker/docker/client"
)

type application struct {
	log          *log.Logger
	orchestrator *containers.Orchestrator
}

func main() {
	port := flag.Int("port", 4000, "port on which the application is exposed")
	flag.Parse()

	cl, err := client.NewEnvClient()
	if err != nil {
		panic(err.Error())
	}

	app := &application{
		log:          log.New(os.Stdout, "", log.Ldate),
		orchestrator: containers.NewOrchestrator(cl),
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: app.routes(),
	}
	log.Printf("Listening on port %d", *port)
	log.Fatal(srv.ListenAndServe())
}
