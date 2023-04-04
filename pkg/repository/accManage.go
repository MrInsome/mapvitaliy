package repository

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	"fmt"
)

func (r *Repository) AddAccount(account types.Account) {
	r.accounts[account.AccountID] = account
}
func (r *Repository) DelAccount(account types.Account) {
	delete(r.accounts, account.AccountID)
}
func (r *Repository) AddAuthData(accountID int) {
	data := types.DataToAccess{
		ClientID:     r.accounts[accountID].Integrations[0].ClientID,
		ClientSecret: r.accounts[accountID].Integrations[0].SecretKey,
		GrantType:    "authorization_code",
		Code:         r.accounts[accountID].Integrations[0].AuthenticationCode,
		RedirectUri:  r.accounts[accountID].Integrations[0].RedirectURL,
	}
	r.data[accountID] = data
}

func (r *Repository) AuthData(accountID int) types.DataToAccess {
	return r.data[accountID]
}

func (r *Repository) AddIntegration(accountID int, integration types.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	account.Integrations = append(account.Integrations, integration)
	r.accounts[accountID] = account
}
func (r *Repository) DelIntegration(accountID int, integration types.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	for i, el := range account.Integrations {
		if el == integration {
			account.Integrations[i] = account.Integrations[len(account.Integrations)-1]
			account.Integrations[len(account.Integrations)-1] = types.Integration{}
			account.Integrations = account.Integrations[:len(account.Integrations)-1]
		}
	}
	r.accounts[accountID] = account
}

func (r *Repository) UpdateIntegration(accountID int, integration types.Integration, replaced types.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	for i, el := range account.Integrations {
		if el == integration {
			account.Integrations[i] = replaced
		}
	}
	r.accounts[accountID] = account

}

func (r *Repository) GetAccountIntegrations(accountID int) []types.Integration {
	account, ok := r.accounts[accountID]
	if !ok {
		return nil
	}
	return account.Integrations
}

func (r *Repository) GetAllAccounts() ([]types.Account, error) {
	accounts := make([]types.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}
	if len(accounts) == 0 {
		return accounts, fmt.Errorf("В базе отсутствуют данные об аккаунтах")
	}
	return accounts, nil
}

func (r *Repository) GetAccount(accountID int) (types.Account, error) {
	if r.accounts[accountID].AccountID == 0 {
		return types.Account{}, fmt.Errorf("aккаунт %d не найден в нашей системе", accountID)
	}
	return r.accounts[accountID], nil
}

func (r *Repository) UnsubscribeAccount(accountID int) error {
	account := r.accounts[accountID]
	if account.AccountID == 0 {
		return fmt.Errorf("аккаунта не существует")
	}
	r.DelAccount(account)
	r.db.Where("account_id = ?", account.AccountID).Delete(&types.Contacts{})
	r.db.Where("account_id = ?", account.AccountID).Delete(&types.Integration{})
	r.db.Delete(account)
	return nil
}

func (r *Repository) SetCurrentAccount() error {
	if config.CurrentAccount == 1 {
		accDB, err := r.GetAllAccounts()
		if err != nil {
			return err
		}
		for i, _ := range accDB {
			if accDB[i].Ref != "" {
				config.CurrentAccount = accDB[i].AccountID
				return nil
			}
		}
		return err
	}
	return nil
}
