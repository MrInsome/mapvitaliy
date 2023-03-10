package internal

import (
	"time"
)

type Account struct {
	AccessToken  string
	RefreshToken string
	Expires      time.Time
	AccountID    int
	Integration  Integration
}

type Integration struct {
	SecretKey          string
	ClientID           string
	RedirectURL        string
	AuthenticationCode string
}

type Repository struct {
	accounts map[int]Account
}

func NewRepository() *Repository {
	return &Repository{
		accounts: make(map[int]Account),
	}
}

func (r *Repository) AddAccount(account Account) {
	r.accounts[account.AccountID] = account
}

func (r *Repository) AddIntegration(accountID int, integration Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	account.Integration = integration
	r.accounts[accountID] = account
}

func (r *Repository) GetAccount(accountID int) (Account, bool) {
	account, ok := r.accounts[accountID]
	return account, ok
}

func (r *Repository) GetAllAccounts() []Account {
	accounts := make([]Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}
