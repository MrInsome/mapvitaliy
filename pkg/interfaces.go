package pkg

import (
	"apitraning/internal"
	"apitraning/internal/types"
	"gorm.io/gorm"
)

type AccountRepo interface {
	AddAccount(account internal.Account)
	GetAccount(accountID int) (internal.Account, error)
	DelAccount(account internal.Account)
	GetAccountIntegrations(accountID int) []internal.Integration
	GetAllAccounts() []internal.Account
	GormDB
}

type GormDB interface {
	SynchronizeDB(db *gorm.DB)
	GormOpen()
	DBReturn() *gorm.DB
}

type RefererRepo interface {
	RefererAdd(ref types.Referer)
	RefererGet() types.Referer
}

type IntegrationRepo interface {
	AddIntegration(accountID int, integration internal.Integration)
	GetAccountIntegrations(accountID int) []internal.Integration
	UpdateIntegration(accountID int, integration internal.Integration, replaced internal.Integration)
	DelIntegration(accountID int, integration internal.Integration)
	GormDB
}

type ContactRepo interface {
	ContactsResp(n types.ContactResponce) []internal.Contacts
	GormDB
}

type AuthRepo interface {
	AddAuthData(accountID int)
	AuthData(accountID int) types.DataToAccess
}

type Unsubscribe interface {
	UnsubscribeAccount(db *gorm.DB, accountID int) error
}

type AccountIntegration interface {
	AccountRepo
	IntegrationRepo
	RefererRepo
	ContactRepo
}

type AccountRefer interface {
	AccountRepo
	RefererRepo
	ContactRepo
	AccountAuth
}

type AccountAuth interface {
	AccountRepo
	AuthRepo
	RefererRepo
}
