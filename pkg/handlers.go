package pkg

import (
	. "apitraning/internal"
	"bytes"
	"encoding/json"
	"net/http"
)

func AuthHandler(repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var respToken TokenResponse
		var data DataToAccess
		account := repo.GetAccount(1)
		data.ClientID = account.Integration[0].ClientID
		data.ClientSecret = account.Integration[0].SecretKey
		data.GrantType = "authorization_code"
		data.Code = account.Integration[0].AuthenticationCode
		data.RedirectUri = account.Integration[0].RedirectURL
		repo.AddAuthData(data)
		a, err := json.Marshal(repo.AuthData())
		resp, err := http.Post("https://testakkamocrm.amocrm.ru/oauth2/access_token",
			"application/json", bytes.NewBuffer(a))
		err = json.NewDecoder(resp.Body).Decode(&respToken)
		err = json.NewEncoder(w).Encode(respToken)
		if err != nil {
			return
		}
		account.AccessToken = respToken.AccessToken
		account.RefreshToken = respToken.RefreshToken
		account.Expires = respToken.ExpiresIn
		repo.AddAccount(account)
		w.WriteHeader(http.StatusCreated)
	}
}
func AmoContat(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := repo.GetAccount(1)

		r, err := http.NewRequest("GET", "https://testakkamocrm.amocrm.ru/api/v4/contacts", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Header.Set("Authorization", "Bearer "+account.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var contacts []Contact

		err = json.NewDecoder(resp.Body).Decode(&contacts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		//err = json.NewEncoder(w).Encode(contacts)
		//if err != nil {
		//	return
		//}

	}
}

func AccountsHandler(repo Repo) http.HandlerFunc {
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
			account.Expires = 86400
			repo.AddAccount(account)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.AddAccount(account)
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.DelAccount(account)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AccountIntegrationsHandler(repo Repo) http.HandlerFunc {
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
				http.Error(w, "Аккаунт не существует", http.StatusNotFound)
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
			integration = account.Integration[0]
			repo.AddIntegration(account.AccountID, integration)
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
			integration = account.Integration[0]
			repo.DelIntegration(account.AccountID, integration)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}

	}

}
