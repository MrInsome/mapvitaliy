package pkg

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	amoCRM2 "apitraning/pkg/integrations/amoCRM"
	"apitraning/pkg/integrations/unisender"
	"apitraning/pkg/repository"
	"encoding/json"
	"net/http"
)

func ContactsSync(repo repository.AccountRefer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.CurrentAccount == 1 {
			accDB, err := repo.GetAllAccounts()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(accDB) > 1 {
				config.CurrentAccount = accDB[len(accDB)-1].AccountID
			} else {
				config.CurrentAccount = accDB[0].AccountID
			}
		}
		amoCRM2.ExportAmo(w, repo)
		account1, err := repo.GetAccount(config.CurrentAccount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = unisender.ImportUni(account1.UniKey, repo, w)
		if err != nil {
			http.Error(w, "Ошибка импорта", http.StatusInternalServerError)
		}

	}
}

func AccountsHandler(repo repository.AccountRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			accounts, err := repo.GetAllAccounts()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.NewEncoder(w).Encode(accounts)
			if err != nil {
				return
			}
		case http.MethodPost:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.AddAccount(account)
			repo.DBReturn().Create(account)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.AddAccount(account)
			repo.DBReturn().Updates(account)
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.DelAccount(account)
			repo.DBReturn().Where("account_id = ?", account.AccountID).Delete(&types.Contacts{})
			repo.DBReturn().Where("account_id = ?", account.AccountID).Delete(&types.Integration{})
			repo.DBReturn().Delete(account)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AccountIntegrationsHandler(repo repository.AccountIntegration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			accIntegration := repo.GetAccountIntegrations(account.AccountID)
			if accIntegration == nil {
				http.Error(w, "Аккаунта или интеграции не существует", http.StatusNotFound)
				return
			}
			err := json.NewEncoder(w).Encode(accIntegration)
			if err != nil {
				return
			}
		case http.MethodPost:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var integration types.Integration
			integration = account.Contactss[0]
			repo.AddIntegration(account.AccountID, integration)
			repo.DBReturn().Updates(integration)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(account.Contactss) < 2 {
				http.Error(w, "Недостаточно аргументов у интеграции", http.StatusBadRequest)
				return
			}
			var integration types.Integration
			integration = account.Contactss[0]
			replaced := account.Contactss[1]
			repo.UpdateIntegration(account.AccountID, integration, replaced)
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var integration types.Integration
			integration = account.Contactss[0]
			repo.DelIntegration(account.AccountID, integration)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}

	}

}

func UnsyncContacts(repo repository.ContactsRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unsyncCon, err := repo.GetUnsyncCon()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(unsyncCon)
		if err != nil {
			return
		}
		return
	}
}
