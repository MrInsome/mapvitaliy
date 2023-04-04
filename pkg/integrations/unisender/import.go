package unisender

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	"apitraning/pkg/repository"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func ImportUni(apiKey string, repo repository.AccountContacts, w http.ResponseWriter) error {
	err := repo.SetCurrentAccount()
	if err != nil {
		return err
	}
	account, err := repo.GetAccount(config.CurrentAccount)
	if err != nil {
		return err
	}
	var uniresp types.ImportUniResponse
	contacts := account.Contacts
	apiUrl := "https://api.unisender.com/ru/api/importContacts"
	for _, el := range contacts {
		data := url.Values{}
		data.Set("format", "json")
		data.Set("api_key", apiKey)
		data.Set("field_names[0]", "email")
		data.Set("field_names[1]", "Name")
		data.Set("data[0][0]", el.Email)
		data.Set("data[0][1]", el.Name)
		resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
		if err != nil {
			return err
		}
		err = json.NewDecoder(resp.Body).Decode(&uniresp)
		if err != nil {
			return err
		}
		if uniresp.Result.Invalid != 0 {
			unSync := types.UnsyncContacts{ContactID: el.ContactID, Name: el.Name}
			repo.AddUnsyncCon(el.ContactID, unSync)
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}
		err = json.NewEncoder(w).Encode(uniresp)
		if err != nil {
			return err
		}
	}
	return nil
}
