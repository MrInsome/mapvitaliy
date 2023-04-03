package repository

import (
	"apitraning/internal/types"
	"github.com/beanstalkd/go-beanstalk"
	"gorm.io/gorm"
)

type Repository struct {
	accounts     map[int]types.Account
	integrations []types.Integration
	contacts     []types.Contacts
	unSyncCon    []types.UnsyncAccounts
	data         map[int]types.DataToAccess
	referer      types.Referer
	db           *gorm.DB
	conn         *beanstalk.Conn
}
