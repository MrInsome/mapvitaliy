package internal

type Account struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	Expires      int           `json:"expires"`
	AccountID    int           `json:"account_id" gorm:"primaryKey:AccountID"`
	Integration  []Integration `gorm:"foreignKey:AccountID"`
	Contact      []Contacts    `gorm:"foreignKey:AccountID"`
}

type Integration struct {
	AccountID          int    `json:"account_id"`
	SecretKey          string `json:"secret_key" gorm:"primaryKey:SecretKey"`
	ClientID           string `json:"client_id"`
	RedirectURL        string `json:"redirect_url"`
	AuthenticationCode string `json:"authentication_code"`
}
type Referer struct {
	Referer string `json:"referer"`
}

type Contacts struct {
	AccountID int    `json:"account_id"`
	Email     string `json:"email" gorm:"primaryKey:Email"`
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DataToAccess struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
}

type ContactsResponse struct {
	Response struct {
		Contacts []struct {
			Name               string `json:"name"`
			Email              string `json:"email"`
			CustomFieldsValues []struct {
				Values []struct {
					Value string `json:"value"`
				} `json:"values"`
			} `json:"custom_fields_values"`
		} `json:"contacts"`
	} `json:"_embedded"`
}
