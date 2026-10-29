package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func main() {
	portNum, err := net.LookupPort("tcp", "http-alt")
	if err != nil {
		log.Fatalf("could not resolve service name: %v", err)
	}

	addr := fmt.Sprintf(":%d", portNum)

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("starting server on %s", addr)

	err = http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
