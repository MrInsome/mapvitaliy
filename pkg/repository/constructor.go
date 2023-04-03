package repository

import (
	"apitraning/internal/types"
)

func NewRepository() *Repository {
	return &Repository{
		accounts:     make(map[int]types.Account),
		integrations: []types.Integration{},
		contacts:     []types.Contacts{},
		unSyncCon:    make(map[int]types.UnsyncContacts),
		data:         make(map[int]types.DataToAccess),
	}
}
