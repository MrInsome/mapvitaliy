package pkg

import "apitraning/internal"

type Repository struct {
	accounts map[int]internal.Account
}

func NewRepository() *Repository {
	return &Repository{
		accounts: make(map[int]internal.Account),
	}
}

func (r *Repository) AddAccount(account internal.Account) {
	r.accounts[account.AccountID] = account
}

func (r *Repository) AddIntegration(accountID int, integration internal.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	account.Integration[integration.ClientID] = integration
	r.accounts[accountID] = account
}

func (r *Repository) GetAccountIntegrations(accountID int) map[string]internal.Integration {
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
