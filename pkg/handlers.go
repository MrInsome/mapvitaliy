package pkg

import (
	. "apitraning/internal"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func FromAMO(repo Repo, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("client_id") != "" {
				var ref Referer
				account := repo.GetAccount(CurrentAccount)
				params := r.URL.Query()
				if params == nil {
					w.WriteHeader(http.StatusConflict)
					return
				}
				account.Integration[CurrentIntegration].AuthenticationCode = params.Get("code")
				account.Integration[CurrentIntegration].ClientID = params.Get("client_id")
				ref.Referer = params.Get("referer")
				repo.RefererAdd(ref)
				repo.AddAccount(account)
				db.Updates(account)
				GetARTokens(repo, db, w)
			}
		case http.MethodPost:
			if r.URL.Query().Get("unisender_key") != "" {
				UniKey = r.URL.Query().Get("unisender_key")
				fmt.Println(UniKey)
			}
		}
	}
}

func AuthHandler(repo Repo, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			GetARTokens(repo, db, w)
		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}

func AmoContact(repo Repo, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contacts ContactResponce
		ref := repo.RefererGet()
		account := repo.GetAccount(CurrentAccount)
		page := 1
		for {
			r, err := http.NewRequest("GET", "https://"+ref.Referer+"/api/v4/contacts?limit=1&page="+strconv.Itoa(page), nil)
			if err != nil {
				http.Error(w, "Неверный запрос", http.StatusInternalServerError)
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
			err = json.NewDecoder(resp.Body).Decode(&contacts)
			if resp.Body == nil {
				break
			}
			account.Contact = repo.ContactsResp(contacts)
			if err != nil {
				return
			}
			repo.AddAccount(account)
			db.Updates(account)
			if err != nil {
				return
			}
			page++
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
			repo.AddAccount(account)
			db.Create(account)
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
			db.Where("account_id = ?", account.AccountID).Delete(&Contacts{})
			db.Where("account_id = ?", account.AccountID).Delete(&Integration{})
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
			integration = account.Integration[CurrentIntegration]
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
			integration = account.Integration[CurrentIntegration]
			repo.DelIntegration(account.AccountID, integration)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}

	}

}

func AdminAccount(repo Repo, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			var integration1 []Integration
			integration1 = append(integration1, Integration{
				SecretKey:          "SffuYBpoDSTUZi0rtW5LPSoJQVV5086bIMScyeSjGyRx9d25ozlUNv0ubLWrkDHj",
				ClientID:           "",
				RedirectURL:        "https://71b5-173-233-147-68.eu.ngrok.io/vidget",
				AuthenticationCode: ""})
			var contact1 []Contacts
			contact1 = append(contact1, Contacts{Email: "yalublugolang@amoschool.zbs"})
			account1 := Account{
				AccountID:   1,
				Integration: integration1,
				Contact:     contact1,
			}
			repo.AddAccount(account1)
			db.Create(&account1)
		}
	}
}
