package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	portNum, err := net.LookupPort("tcp", "http-alt")
	if err != nil {
		log.Fatalf("could not resolve service name: %v", err)
	}

	addr := fmt.Sprintf(":%d", portNum)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Printf("starting server on %s", addr)

	err = http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
