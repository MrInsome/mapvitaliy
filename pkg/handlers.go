package pkg

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	"apitraning/pkg/amoCRM"
	"apitraning/pkg/integrations/unisender"
	"apitraning/pkg/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func FromAMOVidget(repo repository.AccountRefer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Недопустимый метод", config.StatusNotAl)
			return
		}
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			http.Error(w, "client_id не указан", config.StatusConf)
			return
		}
		integration := []types.Integration{
			{
				SecretKey:          config.SecretKey,
				ClientID:           clientID,
				RedirectURL:        config.RedirectURL,
				AuthenticationCode: r.URL.Query().Get("code"),
			},
		}
		account := types.Account{
			AccountID:   config.AccountID,
			Integration: integration,
			Ref:         r.URL.Query().Get("referer"),
		}
		repo.AddAccount(account)
	}
}

func HandleUnisenKey(repo repository.AccountAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
			return
		}
		currentAcc, err := repo.GetAccount(config.CurrentAccount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(r.FormValue("account_id"))
		if err != nil {
			http.Error(w, "Неизвестный Аккаунт ID", http.StatusBadRequest)
			return
		}
		if id != config.CurrentAccount {
			acc := types.Account{
				AccountID:   id,
				UniKey:      r.FormValue("unisender_key"),
				Integration: currentAcc.Integration,
				Ref:         currentAcc.Ref,
			}
			config.CurrentAccount = id
			repo.AddAccount(acc)
			if err := repo.DBReturn().Create(acc).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := amoCRM.GetARTokens(repo, repo.DBReturn(), w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func AuthHandler(repo repository.AccountAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			err := amoCRM.GetARTokens(repo, repo.DBReturn(), w)
			if err != nil {
				http.Error(w, "Ошибка получения токенов авторизации", http.StatusGone)
				return
			}
		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AmoContact(repo repository.AccountRefer) http.HandlerFunc {
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
		amoCRM.ExportAmo(w, repo)
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
			integration = account.Integration[0]
			repo.AddIntegration(account.AccountID, integration)
			repo.DBReturn().Updates(integration)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(account.Integration) < 2 {
				http.Error(w, "Недостаточно аргументов у интеграции", http.StatusBadRequest)
				return
			}
			var integration types.Integration
			integration = account.Integration[0]
			replaced := account.Integration[1]
			repo.UpdateIntegration(account.AccountID, integration, replaced)
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			var account types.Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var integration types.Integration
			integration = account.Integration[0]
			repo.DelIntegration(account.AccountID, integration)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}

	}

}

func UnsyncContacts(repo repository.AccountRepo) http.HandlerFunc {
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

func WebhookFunc(repo repository.BStalkWH) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Не могу прочитать form webhook", http.StatusBadRequest)
			return
		}
		var contact types.Contacts
		add := r.FormValue("contacts[add][0][id]")
		if add != "" {
			_, err := repo.Put([]byte("AddContact "+add), 1, 10, 100*time.Second)
			fmt.Fprint(w, "Работа добавлена в очередь")
			if err != nil {
				return
			}
		}
		update := r.FormValue("contacts[update][0][id]")
		if update != "" {
			_, err := repo.Put([]byte("UpdateContact "+update), 1, 10, 100*time.Second)
			fmt.Fprint(w, "Работа добавлена в очередь")
			if err != nil {
				return
			}
		}
		del := r.FormValue("contacts[delete][0][id]")
		if del != "" {
			_, err := repo.Put([]byte("DeleteContact "+strconv.Itoa(contact.ContactID)), 1, 10, 100*time.Second)
			fmt.Fprint(w, "Работа добавлена в очередь")
			if err != nil {
				return
			}
		}
		return
	}
}
