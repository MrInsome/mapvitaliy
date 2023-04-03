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
	unSyncCon    map[int]types.UnsyncContacts
	data         map[int]types.DataToAccess
	db           *gorm.DB
	conn         *beanstalk.Conn
}
