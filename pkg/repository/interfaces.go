package repository

import (
	"apitraning/internal/types"
	"gorm.io/gorm"
	"time"
)

type AccountRepo interface {
	AddAccount(account types.Account)
	GetAccount(accountID int) (types.Account, error)
	DelAccount(account types.Account)
	GetAccountIntegrations(accountID int) []types.Integration
	GetAllAccounts() ([]types.Account, error)
	AddUnsyncCon(id int, contact types.UnsyncContacts)
	GetUnsyncCon() ([]types.UnsyncContacts, error)
	SetCurrentAccount() error
	GormDB
}

type GormDB interface {
	SynchronizeDB(db *gorm.DB)
	GormOpen() error
	DBReturn() *gorm.DB
}

type BStalkWH interface {
	NewBeanstalkConn() (*Repository, error)
	Close() error
	Put(body []byte, priority uint32, delay, ttr time.Duration) (uint64, error)
	Delete(id uint64) error
	Reserve(ttr time.Duration) (id uint64, body []byte, err error)
}

type IntegrationRepo interface {
	AddIntegration(accountID int, integration types.Integration)
	GetAccountIntegrations(accountID int) []types.Integration
	UpdateIntegration(accountID int, integration types.Integration, replaced types.Integration)
	DelIntegration(accountID int, integration types.Integration)
	GormDB
}

type ContactRepo interface {
	ParseContactsResponse(n types.ContactResponce) []types.Contacts
	GormDB
}

type AuthRepo interface {
	AddAuthData(accountID int)
	AuthData(accountID int) types.DataToAccess
}

type Unsubscribe interface {
	UnsubscribeAccount(accountID int) error
}

type AccountIntegration interface {
	AccountRepo
	IntegrationRepo
	ContactRepo
}

type AccountRefer interface {
	AccountRepo
	ContactRepo
	AccountAuth
}

type AccountAuth interface {
	AccountRepo
	AuthRepo
}
