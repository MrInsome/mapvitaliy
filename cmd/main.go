package main

import (
	. "apitraning/pkg"
	. "apitraning/pkg/api"
	"log"
	"net/http"
)

func main() {
	repo := NewRepository()
	repo.GormOpen()
	repo, err := repo.NewBeanstalkConn()
	if err != nil {
		log.Fatal(err)
	}
	router := Router(repo)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go OpenGRPC(repo)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
