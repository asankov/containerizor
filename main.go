package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/containers", listContainers)
	mux.HandleFunc("/containers/create", createContainer)

	log.Println("Listening on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Containerizor")
}

func listContainers(w http.ResponseWriter, r *http.Request) {
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

func createContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Now Allowed", 405)
		return
	}

	fmt.Fprintln(w, "Creating a new container...")
}
