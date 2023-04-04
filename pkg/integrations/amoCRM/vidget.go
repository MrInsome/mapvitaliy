package amoCRM

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	"apitraning/pkg/repository"
	"net/http"
	"strconv"
)

func FromAMOVidget(repo repository.AccountRefer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Недопустимый метод", config.StatusNotAl)
			return
		}
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			http.Error(w, "client_id не указан", config.StatusConf)
			return
		}
		integration := []types.Integration{
			{
				SecretKey:          config.SecretKey,
				ClientID:           clientID,
				RedirectURL:        config.RedirectURL,
				AuthenticationCode: r.URL.Query().Get("code"),
			},
		}
		account := types.Account{
			AccountID: config.AccountID,
			Contactss: integration,
			Ref:       r.URL.Query().Get("referer"),
		}
		repo.AddAccount(account)
	}
}

func HandleUnisenKey(repo repository.AccountAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
			return
		}
		currentAcc, err := repo.GetAccount(config.CurrentAccount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(r.FormValue("account_id"))
		if err != nil {
			http.Error(w, "Неизвестный Аккаунт ID", http.StatusBadRequest)
			return
		}
		if id != config.CurrentAccount {
			acc := types.Account{
				AccountID: id,
				UniKey:    r.FormValue("unisender_key"),
				Contactss: currentAcc.Contactss,
				Ref:       currentAcc.Ref,
			}
			config.CurrentAccount = id
			repo.AddAccount(acc)
			if err := repo.DBReturn().Create(acc).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := GetARTokens(repo, repo.DBReturn(), w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
