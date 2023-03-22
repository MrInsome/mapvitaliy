package pkg

import (
	"apitraning/internal"
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
}

type Repository struct {
	accounts map[int]internal.Account
	contacts map[int]internal.Contacts
	data     map[int]internal.DataToAccess
	referer  internal.Referer
}

func NewRepository() *Repository {
	return &Repository{
		accounts: make(map[int]internal.Account),
		contacts: make(map[int]internal.Contacts),
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
