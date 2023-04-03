package repository

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (r *Repository) DBReturn() *gorm.DB {
	return r.db
}

func (r *Repository) GormOpen() error {
	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Ошибка соединения с базой данных: %w", err)
	}
	r.db = db
	err = r.db.AutoMigrate(&types.Account{}, &types.Integration{}, &types.Contacts{})
	if err != nil {
		return fmt.Errorf("Ошибка миграции базы данных: %w", err)
	}
	r.SynchronizeDB(r.db)
	return nil
}

func (r *Repository) SynchronizeDB(db *gorm.DB) {
	var account []types.Account
	var contacts []types.Contacts
	var integrations []types.Integration
	db.Find(&account)
	for i, el := range account {
		db.Where("account_id = ?", el.AccountID).Find(&integrations)
		account[i].Integration = integrations
		db.Where("account_id = ?", el.AccountID).Find(&contacts)
		account[i].Contacts = contacts
		r.AddAccount(account[i])
	}
	db.Find(&contacts)
	db.Find(&integrations)
	r.contacts = contacts
	r.integrations = integrations
}
