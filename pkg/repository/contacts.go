package repository

import (
	"apitraning/internal/types"
	"fmt"
	"net/mail"
)

func (r *Repository) ParseContactsResponse(parseModel types.ContactResponce) []types.Contacts {
	for _, v := range parseModel.Response.Contacts {
		id := v.ID
		name := v.Name
		customFields := v.EmailValues
		for _, cf := range customFields {
			if cf.FieldCode == "EMAIL" {
				if cf.Values[0].Value != "" {
					_, err := mail.ParseAddress(cf.Values[0].Value)
					if err == nil {
						r.contacts = append(r.contacts, types.Contacts{Name: name, ContactID: id, Email: cf.Values[0].Value})
					}
				}
			}
		}
	}
	return r.contacts
}

func (r *Repository) AddUnsyncCon(id int, contact types.UnsyncContacts) {
	r.unSyncCon[id] = contact
}
func (r *Repository) GetUnsyncCon() ([]types.UnsyncContacts, error) {
	unsyncContacts := make([]types.UnsyncContacts, 0, len(r.accounts))
	for _, account := range r.unSyncCon {
		unsyncContacts = append(unsyncContacts, account)
	}
	if len(unsyncContacts) == 0 {
		return unsyncContacts, fmt.Errorf("В базе отсутствуют данные о несихронизированных контактах")
	}
	return unsyncContacts, nil
}
