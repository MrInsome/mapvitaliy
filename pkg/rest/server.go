package rest

import (
	"apitraning/pkg/repository"
	"log"
	"net/http"
)

func StartRESTServer(repo *repository.Repository) {
	router := Router(repo)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
