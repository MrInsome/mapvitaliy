package pkg

import (
	"apitraning/internal"
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

type Repo interface {
	AddAccount(account internal.Account)
	DelAccount(account internal.Account)
	GetAccount(accountID int) internal.Account
	RefererAdd(ref internal.Referer)
	RefererGet() internal.Referer
	AddAuthData(accountID int)
	AuthData(accountID int) internal.DataToAccess
	AddIntegration(accountID int, integration internal.Integration)
	DelIntegration(accountID int, integration internal.Integration)
	UpdateIntegration(accountID int, integration internal.Integration, replaced internal.Integration)
	GetAccountIntegrations(accountID int) []internal.Integration
	GetAllAccounts() []internal.Account
	ContactsResp(n internal.ContactResponce) []internal.Contacts
}

type Repository struct {
	accounts map[int]internal.Account
	contacts []internal.Contacts
	data     map[int]internal.DataToAccess
	referer  internal.Referer
}

func NewRepository() *Repository {
	return &Repository{
		accounts: make(map[int]internal.Account),
		contacts: []internal.Contacts{},
		data:     make(map[int]internal.DataToAccess),
		referer:  internal.Referer{},
	}
}
func (r *Repository) RefererAdd(ref internal.Referer) {
	r.referer = ref
}
func (r *Repository) RefererGet() internal.Referer {
	return r.referer
}
func (r *Repository) AddAccount(account internal.Account) {
	r.accounts[account.AccountID] = account
}
func (r *Repository) DelAccount(account internal.Account) {
	delete(r.accounts, account.AccountID)
}
func (r *Repository) AddAuthData(accountID int) {
	var data internal.DataToAccess
	data.ClientID = r.accounts[accountID].Integration[0].ClientID
	data.ClientSecret = r.accounts[accountID].Integration[0].SecretKey
	data.GrantType = "authorization_code"
	data.Code = r.accounts[accountID].Integration[0].AuthenticationCode
	data.RedirectUri = r.accounts[accountID].Integration[0].RedirectURL
	r.data[accountID] = data
}

func (r *Repository) AuthData(accountID int) internal.DataToAccess {
	return r.data[accountID]
}

func (r *Repository) AddIntegration(accountID int, integration internal.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	account.Integration = append(account.Integration, integration)
	r.accounts[accountID] = account
}
func (r *Repository) DelIntegration(accountID int, integration internal.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	for i, el := range account.Integration {
		if el == integration {
			account.Integration[i] = account.Integration[len(account.Integration)-1]
			account.Integration[len(account.Integration)-1] = internal.Integration{}
			account.Integration = account.Integration[:len(account.Integration)-1]
		}
	}
	r.accounts[accountID] = account
}

func (r *Repository) UpdateIntegration(accountID int, integration internal.Integration, replaced internal.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	for i, el := range account.Integration {
		if el == integration {
			account.Integration[i] = replaced
		}
	}
	r.accounts[accountID] = account

}

func (r *Repository) GetAccountIntegrations(accountID int) []internal.Integration {
	account, ok := r.accounts[accountID]
	if !ok {
		return nil
	}
	return account.Integration
}

func (r *Repository) GetAllAccounts() []internal.Account {
	accounts := make([]internal.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}

func (r *Repository) GetAccount(accountID int) internal.Account {
	return r.accounts[accountID]
}
func (r *Repository) ContactsResp(n internal.ContactResponce) []internal.Contacts {
	for _, v := range n.Response.Contacts {
		name := v.Name
		customFields := v.EmailValues
		for _, cf := range customFields {
			if cf.FieldCode == "EMAIL" {
				if cf.Values[0].Value != "" {
					r.contacts = append(r.contacts, internal.Contacts{Name: name, Email: cf.Values[0].Value})
				}
			}
		}
	}
	return r.contacts
}

func GetARTokens(repo Repo, db *gorm.DB, w http.ResponseWriter) {
	account := repo.GetAccount(internal.CurrentAccount)
	ref := repo.RefererGet()
	var respToken internal.TokenResponse
	repo.AddAuthData(internal.CurrentAccount)
	a, err := json.Marshal(repo.AuthData(internal.CurrentAccount))
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
	db.Updates(account)
}

func Router(repo Repo, db *gorm.DB) *http.ServeMux {
	AdminAlphaTest := AdminAccount(repo, db)
	Handler := AccountsHandler(repo, db)
	IntegrationHandler := AccountIntegrationsHandler(repo, db)
	Auth := AuthHandler(repo, db)
	RequestHandler := AmoContact(repo, db)
	GetFromIntegration := FromAMO(repo, db)

	router := http.NewServeMux()

	router.Handle("/vidget", GetFromIntegration)
	router.Handle("/accounts", Handler)
	router.Handle("/access_token", Auth)
	router.Handle("/request", RequestHandler)
	router.Handle("/accounts/integrations", IntegrationHandler)
	router.Handle("/start", AdminAlphaTest)
	return router
}
