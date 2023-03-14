package internal

import (
	"time"
)

type Account struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	Expires      time.Time     `json:"expires"`
	AccountID    int           `json:"account_id"`
	Integration  []Integration `json:"integration"`
}

type Integration struct {
	SecretKey          string `json:"secret_key"`
	ClientID           string `json:"client_id"`
	RedirectURL        string `json:"redirect_url"`
	AuthenticationCode string `json:"authentication_code"`
}
