package repository

import (
	"apitraning/internal/types"
)

func NewRepository() *Repository {
	return &Repository{
		accounts:     make(map[int]types.Account),
		integrations: []types.Integration{},
		contacts:     []types.Contacts{},
		unSyncCon:    []types.UnsyncAccounts{},
		data:         make(map[int]types.DataToAccess),
		referer:      types.Referer{},
	}
}
