package pkg

import (
	. "apitraning/internal"
	"encoding/json"
	"net/http"
	"time"
)

func AccountsHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			accounts := repo.GetAllAccounts()
			err := json.NewEncoder(w).Encode(accounts)
			if err != nil {
				return
			}
		case http.MethodPost:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			account.Expires = time.Now().Add(time.Hour * 24 * 30)
			repo.AddAccount(account)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AccountIntegrationsHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			accIntegration := repo.GetAccountIntegrations(account.AccountID)
			err := json.NewEncoder(w).Encode(accIntegration)
			if err != nil {
				return
			}
		case http.MethodPost:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			//repo.AddIntegration(account.AccountID, account.Integration[])
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}

	}

}
