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
		SecretKey:          "wFNOx9rnLYqIHm7brU19gN7u7k8vmaz2zQycAUGP2X8xVHNuJrgMd8MdfMYEZLEY",
		ClientID:           "1c65fcfb-11d0-4bd6-931f-28b4d89892cf",
		RedirectURL:        "randomurl/redirect",
		AuthenticationCode: "def50200b24579a4dfdc84957826722e78a4338fa53e9637410888b3983a5be99a49a6ea9d5ee1c98e99624daa4f8d603e17f2c95ae81abc1814e521ab5df01f87c1c80a9fb6c47f82a5c5d096840902782aa89a6d253bbf4659542e4508111c9bde7727c9a5544db63d81999e5eca4d949a84b81c1b9c34b9500b767178a18dfa3fe43ea62fcf1a135273480ae910321ba566e3383060487730ecf13acb18400a3cf0298d9ea0d57a860ba53cd02165fb14502145b951865554ad6bc09ab2a1125505c1863c99044a2d200a564683a66fc2c380551b8b675e7d64b9eb036ced3665004b7416f8e2a791ea430da4faedff9e8d669c9e03b4314c02bc21dfc20663d507195b86960a1c398ce7b64d2769d758e2d6f25d09b63298e79e5ac64f5724f857bfeeb0bda1449b9d773f02af82f460138e12a41f5f338d88e35d2a8a2b6cceaa349e32c3a4e8b41b9ba2e07401a71793d1c50e570e827036908c89f8bd72c23d38809f6db2be26dade428011ee71fe1d9debf2a372724aee8b14c5dc5e9e7c47c6a1332d84a9e33d9fab219ae785749945556b1e29bfa3643b1d51ed7f1cb2ddb71c84d29d2f037148c8b88205a7a66940578fc7b40f34f0c9a0537ae67b440494e6ef0c615cb9cc500494cc186d68b09ff22d0c195cc08275153ff8c194d7ae559fe97582"})
	account1 := Account{
		AccessToken:  "access_token_1",
		RefreshToken: "refresh_token_1",
		Expires:      time.Now().Add(time.Hour * 24),
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
