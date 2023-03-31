package pkg

import (
	. "apitraning/internal"
	"apitraning/internal/config"
	"apitraning/internal/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func FromAMOVidget(repo AccountRefer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("client_id") != "" {
				if config.CurrentAccount == 1 {
					var ref types.Referer
					var account Account
					params := r.URL.Query()
					if params == nil {
						w.WriteHeader(http.StatusConflict)
						return
					}
					var integration []Integration
					integration = append(integration, Integration{
						SecretKey:          config.SecretKey,
						ClientID:           params.Get("client_id"),
						RedirectURL:        config.RedirectURL,
						AuthenticationCode: params.Get("code")})
					account = Account{
						AccountID:   1,
						Integration: integration,
					}
					ref.Referer = params.Get("referer")
					repo.RefererAdd(ref)
					repo.AddAccount(account)
				}
			}
		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func UnisenKey(repo AccountAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var blinkAcc Account
			blinkAcc, _ = repo.GetAccount(config.CurrentAccount)
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := strconv.Atoi(r.FormValue("account_id"))
			if id != config.CurrentAccount {
				account := Account{
					AccountID:   id,
					UniKey:      r.FormValue("unisender_key"),
					Integration: blinkAcc.Integration,
				}
				config.CurrentAccount = id
				repo.AddAccount(account)
				repo.DBReturn().Create(account)
				err = GetARTokens(repo, repo.DBReturn(), w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}

	}

}
func AuthHandler(repo AccountAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			err := GetARTokens(repo, repo.DBReturn(), w)
			if err != nil {
				http.Error(w, "Ошибка получения токенов авторизации", http.StatusGone)
				return
			}
		case http.MethodPut:
			var ca types.CurrentAcc
			if err := json.NewDecoder(r.Body).Decode(&ca); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			config.CurrentAccount = ca.Current
		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AmoContact(repo AccountRefer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ExportAmo(w, repo)
		account1, err := repo.GetAccount(config.CurrentAccount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = ImportUni(account1.UniKey, repo, w)
		if err != nil {
			http.Error(w, "Ошибка импорта", http.StatusInternalServerError)
		}

	}
}

func AccountsHandler(repo AccountRepo) http.HandlerFunc {
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
			repo.AddAccount(account)
			repo.DBReturn().Create(account)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.AddAccount(account)
			repo.DBReturn().Updates(account)
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.DelAccount(account)
			repo.DBReturn().Where("account_id = ?", account.AccountID).Delete(&Contacts{})
			repo.DBReturn().Where("account_id = ?", account.AccountID).Delete(&Integration{})
			repo.DBReturn().Delete(account)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AccountIntegrationsHandler(repo AccountIntegration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var account Account
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
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var integration Integration
			integration = account.Integration[config.CurrentIntegration]
			repo.AddIntegration(account.AccountID, integration)
			repo.DBReturn().Updates(integration)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(account.Integration) < 2 {
				http.Error(w, "Недостаточно аргументов у интеграции", http.StatusBadRequest)
				return
			}
			var integration Integration
			integration = account.Integration[0]
			replaced := account.Integration[1]
			repo.UpdateIntegration(account.AccountID, integration, replaced)
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var integration Integration
			integration = account.Integration[config.CurrentIntegration]
			repo.DelIntegration(account.AccountID, integration)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}

	}

}

func WebhookFunc(repo BStalkWH) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for a := 0; a < 10; a++ {
			id, body, _ := repo.Reserve(5 * time.Second)
			idjb, _ := repo.Put([]byte("Job +"), 1, 10, 100*time.Second)
			repo.Delete(id)
			fmt.Fprintf(w, string(body)+"\n")
			fmt.Fprintf(w, strconv.Itoa(int(id))+" "+strconv.Itoa(int(idjb))+"\n")
		}
		id, body, _ := repo.Reserve(5 * time.Second)
		fmt.Fprintf(w, string(body)+"\n"+strconv.Itoa(int(id)))
		repo.Close()
	}
}
