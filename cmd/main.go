package main

import (
	. "apitraning/pkg/api"
	"apitraning/pkg/repository"
	"apitraning/pkg/rest"
	"log"
	"net/http"
)

func main() {
	repo := repository.NewRepository()
	repo.GormOpen()
	repo, err := repo.NewBeanstalkConn()
	if err != nil {
		log.Fatal(err)
	}
	router := rest.Router(repo)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go OpenGRPC(repo)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
