package amoCRM

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	"apitraning/pkg/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

func ExportAmo(w http.ResponseWriter, repo repository.AccountRefer) {
	err := repo.SetCurrentAccount()
	var contacts types.ContactResponce
	account, err := repo.GetAccount(config.CurrentAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ref := account.Ref
	page := 1
	for {
		r, err := http.NewRequest("GET", "https://"+ref+"/api/v4/contacts?limit=1&page="+strconv.Itoa(page), nil)
		if err != nil {
			http.Error(w, "Неверный запрос", http.StatusInternalServerError)
			return
		}
		r.Header.Set("Authorization", "Bearer "+account.AccessToken)
		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&contacts)
		if resp.Body == nil {
			return
		}
		account.Contacts = repo.ParseContactsResponse(contacts)
		if err != nil {
			return
		}
		repo.AddAccount(account)
		repo.DBReturn().Updates(account)
		if err != nil {
			return
		}
		page++
	}
}
