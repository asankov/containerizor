package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asankov/containerizor/internal/db"

	"github.com/asankov/containerizor/internal/db/postgres"

	_ "github.com/lib/pq"

	"github.com/asankov/containerizor/pkg/containers"
	"github.com/docker/docker/client"
)

type application struct {
	log          *log.Logger
	orchestrator *containers.Orchestrator
	db           db.Database
}

const (
	host     = "localhost"
	dbPort   = 5432
	user     = "antonsankov"
	password = ""
	dbname   = "containerizor"
)

func main() {
	port := flag.Int("port", 4000, "port on which the application is exposed")
	flag.Parse()

	cl, err := client.NewEnvClient()
	if err != nil {
		panic(err.Error())
	}

	// TODO: password, sslmode
	db, err := postgres.New(host, dbPort, user, dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := &application{
		log:          log.New(os.Stdout, "", log.Ldate),
		orchestrator: containers.NewOrchestrator(cl),
		db:           db,
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: app.routes(),
	}
	log.Printf("Listening on port %d", *port)
	log.Fatal(srv.ListenAndServe())
}
