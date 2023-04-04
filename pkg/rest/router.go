package rest

import (
	"apitraning/pkg"
	"apitraning/pkg/integrations/amoCRM"
	"apitraning/pkg/repository"
	"net/http"
)

func Router(repo *repository.Repository) *http.ServeMux {
	AccountsHandler := pkg.AccountsHandler(repo)
	IntegrationsHandler := pkg.AccountIntegrationsHandler(repo)
	AuthHandler := amoCRM.AuthHandler(repo)
	ContactsSync := pkg.ContactsSync(repo)
	FromAmoVidget := amoCRM.FromAMOVidget(repo)
	UnisenKey := amoCRM.HandleUnisenKey(repo)
	WebhookProd := amoCRM.WebhookProducer(repo)
	WebhookWork := amoCRM.WebhookWorker(repo)
	Unsync := pkg.UnsyncContacts(repo)
	Sync := pkg.SyncContacts(repo)
	AccountHasUni := pkg.UnisenderCheck(repo)

	router := http.NewServeMux()
	router.Handle("/vidget", FromAmoVidget)
	router.Handle("/vidget/unisender", UnisenKey)
	router.Handle("/accounts", AccountsHandler)
	router.Handle("/accounts/integrations", IntegrationsHandler)
	router.Handle("/accounts/unsync", Unsync)
	router.Handle("/accounts/sync", Sync)
	router.Handle("/access_token", AuthHandler)
	router.Handle("/sync", ContactsSync)
	router.Handle("/webhook", WebhookProd)
	router.Handle("/webhookwork", WebhookWork)
	router.Handle("/acchasuni", AccountHasUni)
	return router
}
