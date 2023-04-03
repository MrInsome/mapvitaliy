package unisender

import (
	"apitraning/internal/config"
	"apitraning/internal/types"
	"apitraning/pkg/repository"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ImportUni(apiKey string, repo repository.AccountRepo, w http.ResponseWriter) error {
	account, err := repo.GetAccount(config.CurrentAccount)
	if err != nil {
		return err
	}
	var uniresp types.ImportUniResponse
	contacts := account.Contacts
	apiUrl := "https://api.unisender.com/ru/api/importContacts"
	data := url.Values{}
	data.Set("format", "json")
	data.Set("api_key", apiKey)
	data.Set("field_names[0]", "email")
	data.Set("field_names[1]", "Name")

	if len(contacts) < 501 {
		for i, el := range contacts {
			data.Set("data["+strconv.Itoa(i)+"][0]", el.Email)
			data.Set("data["+strconv.Itoa(i)+"][1]", el.Name)
		}
	}

	resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	err = json.NewDecoder(resp.Body).Decode(&uniresp)
	err = json.NewEncoder(w).Encode(uniresp)
	for _, el := range uniresp.Result.Log {
		a, err := strconv.Atoi(el.Index)
		if err != nil {
			return err
		}
		repo.AddSyncCon(a, contacts[a])
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
