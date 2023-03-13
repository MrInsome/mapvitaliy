package main

import (
	. "apitraning/internal"
	. "apitraning/pkg"
	"log"
	"net/http"
	"time"
)

func main() {
	repo := NewRepository()
	//Тестовый аккаунт
	var integration1 []Integration
	integration1 = append(integration1, Integration{
		SecretKey:          "secret_key_1",
		ClientID:           "client_id_1",
		RedirectURL:        "randomurl/redirect",
		AuthenticationCode: "auth_code_1"})
	account1 := Account{
		AccessToken:  "access_token_1",
		RefreshToken: "refresh_token_1",
		Expires:      time.Now().Add(time.Hour * 24 * 30),
		AccountID:    1,
		Integration:  []Integration{},
	}
	repo.AddAccount(account1)
	repo.AddIntegration(account1.AccountID, integration1[0])

	handler := AccountsHandler(repo)
	integrationHandler := AccountIntegrationsHandler(repo)

	router := http.NewServeMux()
	router.Handle("/accounts", handler)
	router.Handle("/accounts/integrations", integrationHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
