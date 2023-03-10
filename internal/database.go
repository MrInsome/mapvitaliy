package internal

import (
	"time"
)

type Account struct {
	AccessToken  string
	RefreshToken string
	Expires      time.Time
	AccountID    int
	Integration  map[string]Integration
}

type Integration struct {
	SecretKey          string
	ClientID           string
	RedirectURL        string
	AuthenticationCode string
}
