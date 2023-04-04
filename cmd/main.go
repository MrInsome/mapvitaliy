package main

import (
	"apitraning/internal/config"
	. "apitraning/pkg/api"
	"apitraning/pkg/repository"
	"apitraning/pkg/rest"
	"log"
)

func main() {
	err := config.CallVarsFromENV()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewRepository()
	err = repo.GormOpen()
	if err != nil {
		log.Fatal(err)
	}
	repo.NewBeanstalkConn()
	rest.StartRESTServer(repo)
	go OpenGRPC(repo)
}
