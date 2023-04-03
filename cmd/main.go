package main

import (
	. "apitraning/pkg/api"
	"apitraning/pkg/repository"
	"apitraning/pkg/rest"
	"log"
)

func main() {
	repo := repository.NewRepository()
	err := repo.GormOpen()
	if err != nil {
		log.Fatal(err)
	}
	repo, err = repo.NewBeanstalkConn()
	if err != nil {
		log.Fatal(err)
	}
	rest.StartRESTServer(repo)
	go OpenGRPC(repo)
}
