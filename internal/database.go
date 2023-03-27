package internal

import "gorm.io/gorm"

type Account struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	Expires      int           `json:"expires"`
	AccountID    int           `json:"account_id" gorm:"primaryKey:AccountID"`
	Integration  []Integration `gorm:"foreignKey:AccountID"`
	Contacts     []Contacts    `gorm:"foreignKey:AccountID"`
}

type Integration struct {
	gorm.Model
	AccountID          int    `json:"account_id"`
	SecretKey          string `json:"secret_key" gorm:"primaryKey:SecretKey"`
	ClientID           string `json:"client_id"`
	RedirectURL        string `json:"redirect_url"`
	AuthenticationCode string `json:"authentication_code"`
}

type Contacts struct {
	gorm.Model
	AccountID int    `json:"account_id"`
	ContactID int    `json:"contactID"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"primaryKey:Email"`
}
