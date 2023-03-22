package main

import (
	. "apitraning/internal"
	. "apitraning/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	repo := NewRepository()
	dsn := "steven:here@tcp(127.0.0.1:3306)/fullstack_api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Невозможно подключится к БД")
	}

	//Тестовый аккаунт
	var integration1 []Integration
	integration1 = append(integration1, Integration{
		SecretKey:          "q1ZDbVvj8nVvuMQiBemRGRfz9jxDxYBdJnAJHR885tGg7wqbHwQ99EOKNk51EQox",
		ClientID:           "1c65fcfb-11d0-4bd6-931f-28b4d89892cf",
		RedirectURL:        "https://b028-173-233-147-68.eu.ngrok.io/",
		AuthenticationCode: "def5020077a662a668d8e24108a7f64ffce1868ed8c5cee21741b6b15bb2da43b587a859501033efced82299f7f3a6e53c8a90d16b9cdaec9b871b1c381dccff4a9d0de7835feb52fd750f8e4241490cb3600cb5dc531cedb7be94ceacb797bcea569b4556187b750bcf1af3e406ef8245da05d3465ce112ef03b3436d956b2a39389408e28c282ab924bd3ab240dc5b20c2f8cbea28166957fc3d0dcecc5ce065a2f05202de53fbb9d9950683a5f76171aba8a374cf43b5f8e53c6c837e9eafcae51dcb1c45ff6b7b4770cefc45c25c530f4ec135a1cb29cb01d2e8e3bf686d6d8d18793f641f8a6cefcf7b4302ceae1d7bb5e00683d23b4059ad3f30af006ee3dc19b25a48ae4593f3d3d45ea1e218d917dd05f0b39576907f1b39fb33d3c39345361bab08895239cd6dbb3443f021fe11033a2a21c683d51ef2f144469f3cf6a2640b7598eca6a9f2645c13187abc54e826fc63497a231c3ffa62f5331c6199cc6d2e9509ddf7a3a24e0ee3c4ee7dacc3d5314c7f7d284e65661254f8c921e013c11d96bec5146db42e22af7ca28e1605425fb05826baca5c574dfc2b086831ce3b1d7ba59799d2400fb9b26a3c00c0793ee4eeeb3274e43bb5c7d3ac28013920f186e5960b4aa0a592b835b9a29549de6056cc0417be830745f92e5f5b22759172166c041c1b"})
	var contact1 []Contacts
	contact1 = append(contact1, Contacts{Email: "yalublugolang@amoschool.zbs"})
	account1 := Account{
		AccountID:   1,
		Integration: integration1,
		Contact:     contact1,
	}
	repo.AddAccount(account1)
	//repo.AddIntegration(account1.AccountID, integration1[0])

	handler := AccountsHandler(repo, db)
	integrationHandler := AccountIntegrationsHandler(repo, db)
	auth := AuthHandler(repo, db)
	requestHandler := AmoContact(repo, db)
	getFromIntegration := GetAmoIntegration(repo, db)

	router := http.NewServeMux()
	router.Handle("/", getFromIntegration)
	router.Handle("/accounts", handler)
	router.Handle("/access_token", auth)
	router.Handle("/request", requestHandler)
	router.Handle("/accounts/integrations", integrationHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err = db.AutoMigrate(&Account{}, &Integration{}, &Contacts{})
	if err != nil {
		panic("Невозможно провести миграцию в БД")
	}
	db.Create(&account1)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
