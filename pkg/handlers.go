package pkg

import (
	. "apitraning/internal"
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

func GetAmoIntegration(repo Repo, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ref Referer
		account := repo.GetAccount(1)
		params := r.URL.Query()
		if params == nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
		account.Integration[0].AuthenticationCode = params.Get("code")
		account.Integration[0].ClientID = params.Get("client_id")
		ref.Referer = params.Get("referer")
		repo.RefererAdd(ref)
		repo.AddAccount(account)
		db.Updates(account)
		w.WriteHeader(http.StatusCreated)
	}
}

func AuthHandler(repo Repo, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var respToken TokenResponse
			ref := repo.RefererGet()
			account := repo.GetAccount(1)
			repo.AddAuthData(1)
			a, err := json.Marshal(repo.AuthData(1))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			resp, err := http.Post("https://"+ref.Referer+"/oauth2/access_token",
				"application/json", bytes.NewBuffer(a))
			err = json.NewDecoder(resp.Body).Decode(&respToken)
			err = json.NewEncoder(w).Encode(respToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			account.AccessToken = respToken.AccessToken
			account.RefreshToken = respToken.RefreshToken
			account.Expires = respToken.ExpiresIn
			repo.AddAccount(account)
			w.WriteHeader(http.StatusCreated)
			db.Updates(account)
		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AmoContact(repo Repo, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contacts []Contacts
		var parsedContacts ContactsResponse
		ref := repo.RefererGet()
		account := repo.GetAccount(1)
		contacts = account.Contact
		r, err := http.NewRequest("GET", "https://"+ref.Referer+"/api/v4/contacts", nil)
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

		err = json.NewDecoder(resp.Body).Decode(&parsedContacts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(parsedContacts.Response.Contacts)
		for _, el := range parsedContacts.Response.Contacts {
			contacts = append(contacts, Contacts{Email: el.Email})
		}
		account.Contact = contacts
		repo.AddAccount(account)
		db.Updates(account)
		if err != nil {
			return
		}
	}
}

func AccountsHandler(repo Repo, db *gorm.DB) http.HandlerFunc {
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
			db.Updates(account)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.AddAccount(account)
			db.Updates(account)
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			var account Account
			if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			repo.DelAccount(account)
			db.Delete(account)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AccountIntegrationsHandler(repo Repo, db *gorm.DB) http.HandlerFunc {
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
			integration = account.Integration[0]
			repo.AddIntegration(account.AccountID, integration)
			db.Updates(integration)
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
