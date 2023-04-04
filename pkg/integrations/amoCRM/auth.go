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
	repo.AddAuthData(config.CurrentAccount)
	auth := repo.AuthData(1)
	if auth.ClientID == "" {
		auth = repo.AuthData(config.CurrentAccount)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	var respToken TokenResponse
	a, err := json.Marshal(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	resp, err := http.Post("https://"+account.Ref+"/oauth2/access_token",
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

func AuthHandler(repo repository.AccountAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			err := GetARTokens(repo, repo.DBReturn(), w)
			if err != nil {
				http.Error(w, "Ошибка получения токенов авторизации", http.StatusGone)
				return
			}
		default:
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		}
	}
}
