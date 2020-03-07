package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/containers", listContainers)
	mux.HandleFunc("/containers/create", createContainer)

	log.Println("Listening on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
