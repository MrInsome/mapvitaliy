package amoCRM

import (
	"apitraning/internal/config"
	"apitraning/pkg/repository"
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GetARTokens(repo repository.AccountAuth, db *gorm.DB, w http.ResponseWriter) error {
	account, err := repo.GetAccount(config.CurrentAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	ref := repo.RefererGet()
	var respToken TokenResponse
	repo.AddAuthData(config.CurrentAccount)
	a, err := json.Marshal(repo.AuthData(config.CurrentAccount))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	resp, err := http.Post("https://"+ref.Referer+"/oauth2/access_token",
		"application/json", bytes.NewBuffer(a))
	err = json.NewDecoder(resp.Body).Decode(&respToken)
	err = json.NewEncoder(w).Encode(respToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	account.AccessToken = respToken.AccessToken
	account.RefreshToken = respToken.RefreshToken
	account.Expires = respToken.ExpiresIn
	repo.AddAccount(account)
	db.Updates(account)
	return nil
}
