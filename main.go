package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Println("Listening on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Containerizor"))
}
