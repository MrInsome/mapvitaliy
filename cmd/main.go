package main

import (
	"log"
	"net/http"
)

func main() {
	repo := NewRepository()
	handler := NewHandler(repo)

	router := http.NewServeMux()
	router.Handle("/accounts", handler)
	router.Handle("/integrations", handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
