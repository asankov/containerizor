package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asankov/containerizor/internal/db"

	"github.com/asankov/containerizor/internal/db/postgres"

	// to register PostreSQL driver
	_ "github.com/lib/pq"

	"github.com/asankov/containerizor/pkg/containers"
	"github.com/docker/docker/client"
)

type server struct {
	log          *log.Logger
	orchestrator *containers.Orchestrator
	db           db.Database
}

func main() {
	port := flag.Int("port", 4000, "port on which the application is exposed")
	dbHost := flag.String("db_host", "localhost", "the address of the database")
	dbPort := flag.Int("db_port", 5432, "the port of the database")
	dbUser := flag.String("db_user", "", "the user of the database")
	dbPass := flag.String("db_pass", "", "the password for the database")
	dbName := flag.String("db_name", "", "the name of the database")
	flag.Parse()

	cl, err := client.NewEnvClient()
	if err != nil {
		panic(err.Error())
	}

	// TODO: password, sslmode
	db, err := postgres.New(*dbHost, *dbPort, *dbUser, *dbName, *dbPass)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := &server{
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
