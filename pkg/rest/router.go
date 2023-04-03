package rest

import (
	"apitraning/pkg"
	"apitraning/pkg/repository"
	"net/http"
)

func Router(repo *repository.Repository) *http.ServeMux {
	Handler := pkg.AccountsHandler(repo)
	IntegrationHandler := pkg.AccountIntegrationsHandler(repo)
	Auth := pkg.AuthHandler(repo)
	RequestHandler := pkg.AmoContact(repo)
	GetFromAmoVidget := pkg.FromAMOVidget(repo)
	FromAmoUniKey := pkg.UnisenKey(repo)
	Webhook := pkg.WebhookFunc(repo)

	router := http.NewServeMux()

	router.Handle("/vidget", GetFromAmoVidget)
	router.Handle("/accounts", Handler)
	router.Handle("/access_token", Auth)
	router.Handle("/request", RequestHandler)
	router.Handle("/accounts/integrations", IntegrationHandler)
	router.Handle("/vidget/unisender", FromAmoUniKey)
	router.Handle("/webhook", Webhook)
	return router
}
