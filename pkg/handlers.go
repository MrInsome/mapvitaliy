package pkg

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	amoCRM2 "apitraning/pkg/integrations/amoCRM"
	"apitraning/pkg/integrations/unisender"
	"apitraning/pkg/repository"
	"encoding/json"
	"fmt"
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
			repo.DBReturn().Where("account_id = ?", account.AccountID).Updates(account)
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
			integration = account.Integrations[0]
			repo.AddIntegration(account.AccountID, integration)
			repo.DBReturn().Updates(integration)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(account.Integrations) < 2 {
				http.Error(w, "Недостаточно аргументов у интеграции", http.StatusBadRequest)
				return
			}
			var integration types.Integration
			integration = account.Integrations[0]
			replaced := account.Integrations[1]
			repo.UpdateIntegration(account.AccountID, integration, replaced)
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var integration types.Integration
			integration = account.Integrations[0]
			repo.DelIntegration(account.AccountID, integration)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}

	}

}

func UnsyncContacts(repo repository.AccountContacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repo.SetCurrentAccount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
func SyncContacts(repo repository.AccountContacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repo.SetCurrentAccount()
		con, err := repo.GetSyncCon()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(con)
		if err != nil {
			return
		}
		return
	}
}

func UnisenderCheck(repo repository.AccountRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repo.SetCurrentAccount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var bl bool
		account, err := repo.GetAccount(config.CurrentAccount)
		if account.UniKey != "" {
			bl = true
			fmt.Fprintf(w, "Аккаунт подключен к Unisender %v", bl)
			return
		} else {
			bl = true
			fmt.Fprintf(w, "Аккаунт не подключен к Unisender %v", bl)
			return
		}
	}
}
