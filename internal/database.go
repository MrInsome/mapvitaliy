package internal

var CurrentAccount = 1
var CurrentIntegration = 0

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
	Name      string `json:"name"`
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

type ContactResponce struct {
	Response struct {
		Contacts []struct {
			//ID          int    `json:"id"`
			Name string `json:"name"`
			//FirstName   string `json:"first_name"`
			//LastName    string `json:"last_name"`
			EmailValues []struct {
				FieldCode string `json:"field_code"`
				Values    []struct {
					Value string `json:"value"`
					//EmailCode string `json:"field_code"`
				} `json:"values"`
			} `json:"custom_fields_values"`
		} `json:"contacts"`
	} `json:"_embedded"`
}
