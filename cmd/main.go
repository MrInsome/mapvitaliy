package main

import (
	. "apitraning/internal"
	. "apitraning/internal/config"
	. "apitraning/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	repo := NewRepository()
	db, err := gorm.Open(mysql.Open(Dsn), &gorm.Config{})
	if err != nil {
		panic("Невозможно подключится к БД")
	}

	router := Router(repo)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err = db.AutoMigrate(&Account{}, &Integration{}, &Contacts{})
	if err != nil {
		panic("Невозможно провести миграцию в БД")
	}
	repo.SynchronizeDB(db)
	//OpenGRPC()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
