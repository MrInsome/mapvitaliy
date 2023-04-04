package types

type Account struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	Expires      int           `json:"expires"`
	UniKey       string        `json:"unisender_key"`
	AccountID    int           `json:"account_id" gorm:"primaryKey:AccountID"`
	Ref          string        `json:"referer"`
	Integrations []Integration `gorm:"foreignKey:AccountID"`
	Contacts     []Contacts    `gorm:"foreignKey:AccountID"`
}

type Integration struct {
	AccountID          int    `json:"account_id"`
	SecretKey          string `json:"secret_key" gorm:"primaryKey:SecretKey"`
	ClientID           string `json:"client_id"`
	RedirectURL        string `json:"redirect_url"`
	AuthenticationCode string `json:"authentication_code"`
}

type Contacts struct {
	AccountID int    `json:"account_id"`
	ContactID int    `json:"contactID"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"primaryKey:Email"`
}

type UnsyncContacts struct {
	ContactID int    `json:"contactID"`
	Name      string `json:"name"`
}
